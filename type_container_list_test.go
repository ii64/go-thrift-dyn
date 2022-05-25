package go_thriftproxy

import (
	"bytes"
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypeContainer_List_Set(t *testing.T) {
	var err error
	is := require.New(t)
	withProtocols(t, defaultTestTProtocols, nil, func(protofactory thrift.TProtocolFactory) {
		withBytesTTransport(t, func(b *bytes.Buffer, trans thrift.TTransport) {
			ctx := context.Background()
			proto := protofactory.GetProtocol(trans)

			data := []string{"hello", "world", "你好", "世界"}

			err = proto.WriteListBegin(ctx, thrift.STRING, len(data))
			is.NoError(err)
			for _, s := range data {
				err = proto.WriteString(ctx, s)
				is.NoError(err)
			}
			err = proto.WriteListEnd(ctx)
			is.NoError(err)

			err = proto.Flush(ctx)
			is.NoError(err)

			var expected = make([]byte, b.Len())
			copy(expected, b.Bytes())
			b.Reset()

			// try to rebuild the same data.
			// thrift.SET and thrift.LIST wire data is exactly the same
			// both are using {write,read}Collection
			// the only difference is on the struct TField type.

			for _, typ := range []TypeContainerImplementer{
				NewTypeContainerList(thrift.LIST, TypeContainerDesc{
					Value: thrift.STRING,
				}, true),
				NewTypeContainerSet(thrift.SET, TypeContainerDesc{
					Value: thrift.STRING,
				}, true),
			} {
				for _, s := range data {
					err = typ.Add(typ.Element().SetValue(s))
					is.NoError(err)
				}
				err = typ.Write(ctx, proto)
				is.NoError(err)

				err = proto.Flush(ctx)
				is.NoError(err)

				var actual = make([]byte, b.Len())
				copy(actual, b.Bytes())
				b.Reset()

				// fmt.Println(actual)

				is.Equalf(expected, actual, "%+#v", typ)
			}
		})
	})
}
