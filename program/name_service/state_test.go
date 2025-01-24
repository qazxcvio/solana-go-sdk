package name_service

import (
	"fmt"
	"testing"

	"github.com/qazxcvio/solana-go-sdk/common"
	"github.com/stretchr/testify/assert"
)

func TestNameRecordHeaderFromData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want NameRecordHeader
		err  error
	}{
		{
			args: args{
				data: []byte{0x3d, 0x53, 0xc2, 0x4b},
			},
			want: NameRecordHeader{},
			err:  fmt.Errorf("data length should bigger than 96"),
		},
		{
			args: args{
				data: []byte{0x3d, 0x53, 0xc2, 0x4b, 0x38, 0x36, 0xe, 0xd3, 0x81, 0x3a, 0x23, 0xdf, 0xb2, 0xdf, 0xd8, 0x20, 0xab, 0x58, 0x21, 0xcb, 0x79, 0x29, 0xa3, 0x8d, 0x2e, 0xaa, 0xb2, 0x52, 0xe8, 0x38, 0x25, 0x95, 0x58, 0x7f, 0x6a, 0x3d, 0xab, 0x65, 0xe7, 0x3e, 0x12, 0xde, 0x67, 0xbc, 0x31, 0x73, 0x2d, 0xa0, 0x4e, 0xea, 0xfb, 0x12, 0x83, 0xdd, 0x21, 0x10, 0x82, 0x5c, 0xcb, 0x1e, 0xdf, 0x79, 0xa2, 0xb0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x52, 0x65, 0x61, 0x63, 0x68, 0x20, 0x6f, 0x75, 0x74, 0x20, 0x74, 0x6f, 0x20, 0x6a, 0x61, 0x63, 0x6b, 0x68, 0x6f, 0x6c, 0x6d, 0x37, 0x37, 0x32, 0x40, 0x67, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6d, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x61, 0x20, 0x73, 0x61, 0x6c, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			want: NameRecordHeader{
				ParentName: common.PublicKeyFromString("58PwtjSDuFHuUkYjH9BYnnQKHfwo9reZhC2zMJv9JPkx"),
				Class:      common.PublicKey{},
				Owner:      common.PublicKeyFromString("6xTZhtNA8aaipc2hHFP616gFvDcvWmYMGsDFHwrsF3m1"),
				Data:       []byte{0x52, 0x65, 0x61, 0x63, 0x68, 0x20, 0x6f, 0x75, 0x74, 0x20, 0x74, 0x6f, 0x20, 0x6a, 0x61, 0x63, 0x6b, 0x68, 0x6f, 0x6c, 0x6d, 0x37, 0x37, 0x32, 0x40, 0x67, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6d, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x61, 0x20, 0x73, 0x61, 0x6c, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NameRecordHeaderFromData(tt.args.data)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
