package thrift_dyn

import (
	"bytes"
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEnsureTTypeStringInterface(t *testing.T) {
	withBytesTTransport(t, func(b *bytes.Buffer, trans thrift.TTransport) {
		withProtocols(t, defaultTestTProtocols, &thrift.TConfiguration{}, func(protofactory thrift.TProtocolFactory) {
			var err error
			ctx := context.Background()
			prot := protofactory.GetProtocol(trans)
			var expected any
			expected = "hello"

			spec := TDataSpec{
				Type:     thrift.STRING,
				Required: true,
				Protocol: prot,
			}
			err = WriteDataGeneric(ctx, spec, expected)
			require.NoError(t, err)
			err = trans.Flush(ctx)
			require.NoError(t, err)

			bs := b.Bytes()

			var actual any
			err = ReadDataGeneric(ctx, spec, &actual)
			require.NoError(t, err)

			// STRING will be decoded as []byte
			// if given interface is not initialized with "" (string)
			require.Equal(t, []byte(expected.(string)), actual, "ensure []byte")

			b.Reset()
			_, err = trans.Write(bs)
			require.NoError(t, err)
			err = trans.Flush(ctx)
			require.NoError(t, err)

			actual = ""
			err = ReadDataGeneric(ctx, spec, &actual)
			require.NoError(t, err)
			require.Equal(t, expected, actual, "ensure string")
		})
	})
}

func TestEnsureTTypeStringGeneric(t *testing.T) {
	withBytesTTransport(t, func(b *bytes.Buffer, trans thrift.TTransport) {
		withProtocols(t, defaultTestTProtocols, &thrift.TConfiguration{}, func(protofactory thrift.TProtocolFactory) {
			var err error
			ctx := context.Background()
			prot := protofactory.GetProtocol(trans)

			spec := TDataSpec{
				Type:     thrift.STRING,
				Required: true,
				Protocol: prot,
			}

			var expected = "hello"
			err = WriteData[string](ctx, TData[string]{
				TDataSpec: spec,
				Value:     &expected,
			})
			require.NoError(t, err)
			err = trans.Flush(ctx)
			require.NoError(t, err)

			bs := b.Bytes()

			var actual []byte
			err = ReadData[[]byte](ctx, TData[[]byte]{
				TDataSpec: spec,
				Value:     &actual,
			})
			require.NoError(t, err)
			require.Equal(t, []byte(expected), actual, "ensure []byte")

			b.Reset()
			_, err = trans.Write(bs)
			require.NoError(t, err)
			err = trans.Flush(ctx)
			require.NoError(t, err)

			var actual2 string
			err = ReadData[string](ctx, TData[string]{
				TDataSpec: spec,
				Value:     &actual2,
			})
			require.NoError(t, err)
			require.Equal(t, expected, actual2, "ensure string")
		})
	})
}
