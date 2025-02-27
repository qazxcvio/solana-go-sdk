package client

import (
	"context"
	"testing"

	"github.com/qazxcvio/solana-go-sdk/internal/client_test"
)

func TestClient_GetSlotLeaders(t *testing.T) {
	client_test.TestAll(
		t,
		[]client_test.Param{
			{
				RequestBody:  `{"jsonrpc":"2.0", "id":1, "method": "getSlotLeaders", "params": [264431722, 10]}`,
				ResponseBody: `{"jsonrpc":"2.0","result":["DUND26mEDfFeaPsVof3YvbXDRvpuQX7HMUJrLgEWzYw4","DUND26mEDfFeaPsVof3YvbXDRvpuQX7HMUJrLgEWzYw4","3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt","3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt","3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt","3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt","BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH","BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH","BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH","BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH"],"id":1}`,
				F: func(url string) (any, error) {
					c := NewClient(url)
					return c.GetSlotLeaders(context.TODO(), 264431722, 10)
				},
				ExpectedValue: []string{
					"DUND26mEDfFeaPsVof3YvbXDRvpuQX7HMUJrLgEWzYw4",
					"DUND26mEDfFeaPsVof3YvbXDRvpuQX7HMUJrLgEWzYw4",
					"3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt",
					"3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt",
					"3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt",
					"3JotfSFPaod4KVK7nj7ULvcq5PjUBdZNVGracNkJNhrt",
					"BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH",
					"BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH",
					"BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH",
					"BXAxLMMMUNYfC1z166VjWHR3WjTmqzLxB837o5ghmRtH",
				},
				ExpectedError: nil,
			},
		},
	)
}
