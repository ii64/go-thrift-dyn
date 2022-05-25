package thrift_dyn

import (
	"bytes"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/ii64/go-thrift-dyn/internal/test/base"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncoder(t *testing.T) {
	var err error
	protos := []struct {
		N        int
		prot     ProtocolType
		expected []byte
	}{
		{2, ProtocolType_Compact, []byte{0x18, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x36, 0xfc, 0xab, 0x6, 0x57, 0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xf0, 0x3f, 0x0}},
		{2, ProtocolType_Binary, []byte{0xb, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xca, 0xfe, 0x4, 0x0, 0x9, 0x3f, 0xf0, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd, 0x0}},
	}
	for _, tc := range protos {
		enc := NewEncoder(ProtocolFactory(tc.prot, &thrift.TConfiguration{}))
		var m = base.Model{
			Abc: "hello",
			Sd:  0xcafe,
			F64: 1.05,
		}
		var bb []byte
		t.Run("encode", func(t *testing.T) {
			for i := 0; i < tc.N; i++ {
				bb, err = enc.Encode(&m)
				require.NoError(t, err)
				require.Equal(t, tc.expected, bb)
			}
		})
		t.Run("encode-to", func(t *testing.T) {
			var b [4096]byte
			var n int
			for i := 0; i < tc.N; i++ {
				n, err = enc.EncodeTo(b[:], &m)
				require.NoError(t, err)
				require.Equal(t, tc.expected, b[:n])
			}
		})
		t.Run("write-to", func(t *testing.T) {
			var b bytes.Buffer
			var n int64
			for i := 0; i < tc.N; i++ {
				b.Reset()
				n, err = enc.WriteTo(&b, &m)
				require.NoError(t, err)
				require.Equal(t, tc.expected, b.Bytes())
				require.Equal(t, len(tc.expected), int(n))
			}
		})
	}
}
