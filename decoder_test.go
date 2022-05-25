package thrift_dyn

import (
	"bytes"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/ii64/go-thrift-dyn/internal/test/base"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecode(t *testing.T) {
	var err error
	expected := base.Model{
		Abc: "hello",
		Sd:  0xcafe,
		F64: 1.05,
	}
	protos := []struct {
		N    int
		prot ProtocolType
		bb   []byte
	}{
		{2, ProtocolType_Compact, []byte{0x18, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x36, 0xfc, 0xab, 0x6, 0x57, 0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xf0, 0x3f, 0x0}},
		{2, ProtocolType_Binary, []byte{0xb, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0xa, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xca, 0xfe, 0x4, 0x0, 0x9, 0x3f, 0xf0, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd, 0x0}},
	}
	for _, tc := range protos {
		dec := NewDecoder(ProtocolFactory(tc.prot, &thrift.TConfiguration{}))
		t.Run("decode", func(t *testing.T) {
			for i := 0; i < tc.N; i++ {
				var actual base.Model
				err = dec.Decode(tc.bb, &actual)
				require.NoError(t, err)
				require.Equal(t, expected, actual)
			}
		})
		t.Run("ReadFrom", func(t *testing.T) {
			for i := 0; i < tc.N; i++ {
				var actual base.Model
				err = dec.ReadFrom(bytes.NewReader(tc.bb), &actual)
				require.NoError(t, err)
				require.Equal(t, expected, actual)
			}
		})

	}
}
