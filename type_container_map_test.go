package go_thriftproxy

import (
	"bytes"
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypeContainer_Map(t *testing.T) {
	var err error
	is := require.New(t)
	withProtocols(t, defaultTestTProtocols, nil, func(protofactory thrift.TProtocolFactory) {
		withBytesTTransport(t, func(b *bytes.Buffer, trans thrift.TTransport) {
			ctx := context.Background()
			proto := protofactory.GetProtocol(trans)

			order := []int64{1, 3}
			data := map[int64]int64{
				1: 2,
				3: 4,
			}

			err = proto.WriteMapBegin(ctx, thrift.I64, thrift.I64, len(data))
			is.NoError(err)
			for _, k := range order {
				v := data[k]
				err = proto.WriteI64(ctx, k)
				is.NoError(err)
				err = proto.WriteI64(ctx, v)
				is.NoError(err)
			}
			err = proto.WriteMapEnd(ctx)
			is.NoError(err)

			err = proto.Flush(ctx)
			is.NoError(err)

			var expected = make([]byte, b.Len())
			copy(expected, b.Bytes())

			// fmt.Println(expected)
			b.Reset()

			// try to rebuild the same data.

			typ := NewTypeContainerMap(thrift.MAP, TypeContainerDesc{thrift.I64, thrift.I64}, true)
			typ.SetSize(len(data))
			for _, k := range order {
				v := data[k]
				err = typ.Add(typ.Element().SetKey(k).SetValue(v))
				is.NoError(err)
			}
			err = typ.Write(ctx, proto)
			is.NoError(err)

			err = proto.Flush(ctx)
			is.NoError(err)

			var actual = make([]byte, b.Len())
			copy(actual, b.Bytes())

			is.Equal(expected, actual, "%+#v", typ)
		})
	})
}
