package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"net/http"
)

var jsoni = jsoniter.Config{
	EscapeHTML:             false, // 禁用 HTML 转义
	SortMapKeys:            false, // 禁止排序 Map 键
	ValidateJsonRawMessage: true,  // 校验 RawMessage
	UseNumber:              true,  // 避免 float64 转换
	DisallowUnknownFields:  true,  // 禁止未知字段
	CaseSensitive:          true,  // 区分大小写（减少匹配开销）

	//EscapeHTML:             false,
	//SortMapKeys:            false,
	//ValidateJsonRawMessage: true,
	//UseNumber:              true, // 避免 float64 转换
	//DisallowUnknownFields:  true, // 禁止未知字段
	//TagKey:                 "json",
}.Froze()

const (
	LocalnetRPCEndpoint = "http://localhost:8899"
	DevnetRPCEndpoint   = "https://api.devnet.solana.com"
	TestnetRPCEndpoint  = "https://api.testnet.solana.com"
	MainnetRPCEndpoint  = "https://api.mainnet-beta.solana.com"
)

type JsonRpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params,omitempty"`
}

type JsonRpcResponse[T any] struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Result  T             `json:"result"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

func (j JsonRpcResponse[T]) GetResult() T {
	return j.Result
}

func (j JsonRpcResponse[T]) GetError() error {
	if j.Error == nil {
		return nil
	}
	return j.Error
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *JsonRpcError) Error() string {
	s, err := json.Marshal(e)
	if err == nil {
		return string(s)
	}

	// ideally, it should never reach here
	return fmt.Sprintf("failed to marshal JsonRpcError, err: %v, code: %v, message: %v, data: %v", err, e.Code, e.Message, e.Data)
}

type ValueWithContext[T any] struct {
	Context Context `json:"context"`
	Value   T       `json:"value"`
}

type RpcClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewRpcClient(endpoint string) RpcClient { return New(WithEndpoint(endpoint)) }

// New applies the given options to the rpc client being created. if no options
// is passed, it defaults to a bare bone http client and solana mainnet
func New(opts ...Option) RpcClient {

	client := &RpcClient{}

	setDefaultOptions(client)

	for _, opt := range opts {
		opt(client)
	}

	return *client
}

// Call will return body of response. if http code beyond 200~300, the error also returns.
func (c *RpcClient) Call(ctx context.Context, params ...any) ([]byte, error) {
	// prepare payload
	j, err := preparePayload(params)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare payload, err: %v", err)
	}

	// prepare request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(j))
	if err != nil {
		return nil, fmt.Errorf("failed to do http.NewRequestWithContext, err: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// do request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request, err: %v", err)
	}
	defer res.Body.Close()

	// parse body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, err: %v", err)
	}

	// check response code
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return body, fmt.Errorf("get status code: %v", res.StatusCode)
	}

	return body, nil
}

func (c *RpcClient) CallStream(ctx context.Context, params ...any) (io.ReadCloser, error) {
	// 准备请求负载
	j, err := preparePayload(params)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare payload, err: %v", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(j))
	if err != nil {
		return nil, fmt.Errorf("failed to create request, err: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed, err: %v", err)
	}

	// 检查状态码
	if res.StatusCode < 200 || res.StatusCode > 300 {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body) // 仅读取错误信息
		return nil, fmt.Errorf("invalid status code: %d, body: %s", res.StatusCode, body)
	}

	// 直接返回 Body 的流式接口
	return res.Body, nil
}

func preparePayload(params []any) ([]byte, error) {
	// prepare payload
	j, err := json.Marshal(JsonRpcRequest{
		JsonRpc: "2.0",
		Id:      1,
		Method:  params[0].(string),
		Params:  params[1:],
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

//func call[T any](c *RpcClient, ctx context.Context, params ...any) (T, error) {
//	var output T
//
//	// rpc call
//	body, err := c.Call(ctx, params...)
//	if err != nil {
//		return output, fmt.Errorf("rpc: call error, err: %v, body: %v", err, string(body))
//	}
//
//	// transfer data
//	err = json.Unmarshal(body, &output)
//	if err != nil {
//		return output, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
//	}
//
//	return output, nil
//}

//func call[T any](c *RpcClient, ctx context.Context, params ...any) (T, error) {
//	var output T
//
//	// rpc call
//	body, err := c.Call(ctx, params...)
//	if err != nil {
//		return output, fmt.Errorf("rpc: call error, err: %v, body: %v", err, string(body))
//	}
//
//	// Use streaming JSON decoding
//	decoder := json.NewDecoder(bytes.NewReader(body))
//
//	// Incrementally decode the JSON response
//	if err := decoder.Decode(&output); err != nil {
//		if err == io.EOF {
//			// No content to decode
//			return output, nil
//		}
//		return output, fmt.Errorf("rpc: failed to json decode body, err: %v", err)
//	}
//	return output, nil
//}

func call[T any](c *RpcClient, ctx context.Context, params ...any) (T, error) {
	var output T

	// 获取流式响应体
	reader, err := c.CallStream(ctx, params...)
	if err != nil {
		return output, fmt.Errorf("rpc call failed: %v", err)
	}
	defer reader.Close()

	// 流式 JSON 解码
	decoder := jsoni.NewDecoder(reader)
	if err := decoder.Decode(&output); err != nil {
		return output, fmt.Errorf("json decode failed: %v", err)
	}

	return output, nil
}
