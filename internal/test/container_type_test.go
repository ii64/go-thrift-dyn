package test

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	th "github.com/ii64/go-thrift-dyn"
	"github.com/ii64/go-thrift-dyn/internal/test/base"
	"github.com/stretchr/testify/require"
	"testing"
	_ "unsafe"
)

//go:noescape
//go:linkname request_writeField44 github.com/ii64/go-thrift-dyn/internal/test/base.(*Request).writeField44
//goland:noinspection GoUnusedParameter
func request_writeField44(bs *base.Request, ctx context.Context, oprot thrift.TProtocol) (err error)

func TestStructField_List_use_Model(t *testing.T) {
	var err error
	protofactory := th.ProtocolFactory(th.ProtocolType_Binary, &thrift.TConfiguration{})
	enc := th.NewEncoder(protofactory)

	trans := thrift.NewTMemoryBufferLen(2 << 11)
	prot := protofactory.GetProtocol(trans)

	var m base.Request
	m.Models = append(m.Models, base.NewModel())
	m.Models = append(m.Models, base.NewModel())

	err = request_writeField44(&m, context.Background(), prot)
	require.NoError(t, err)

	var expected = trans.Bytes()

	typ := th.NewTypeContainerList[*base.Model](th.TypeContainerDesc{Value: thrift.STRUCT}, true)
	typ.Value = m.Models

	f1 := th.NewTField(44, thrift.LIST, "models", true)
	f1.SetValue(typ)
	var actual []byte
	actual, err = enc.Encode(f1)
	require.NoError(t, err)

	fmt.Println(actual)
	require.Equal(t, expected, actual)
}

func TestStructField_List_Map(t *testing.T) {
	var err error
	enc := th.NewEncoder(th.ProtocolFactory(th.ProtocolType_Binary, &thrift.TConfiguration{}))

	var n = base.SimpleListMap{
		ListListI32: [][]int32{
			{1, 2, 3, 4},
		},
		ListI64:  []int64{5, 6, 7, 8},
		I64ByI64: nil,
	}
	var expected []byte
	{
		expected, err = enc.Encode(&n)
		require.NoError(t, err)
	}

	var m th.RPCStruct
	{
		f1 := th.NewTField(10, thrift.LIST, "listListI32", true)
		listListI32 := th.NewTypeContainerList[th.TypeContainerImplementer](th.TypeContainerDesc{Value: thrift.LIST}, true)
		f1.SetValue(listListI32)

		f2 := th.NewTField(555, thrift.LIST, "listI64", true)
		listI64 := th.NewTypeContainerList[int64](th.TypeContainerDesc{Value: thrift.I64}, true)
		f2.SetValue(listI64)

		f3 := th.NewTField(2016, thrift.MAP, "i64ByI64", true)
		f3v := th.NewTypeContainerMap[int64, int64](th.TypeContainerDesc{Key: thrift.I64, Value: thrift.I64}, true)
		f3.SetValue(f3v)

		m.AddField(f1, f2, f3)

		for _, listI32Value := range n.ListListI32 {
			listI32 := th.NewTypeContainerList[int32](th.TypeContainerDesc{Value: thrift.I32}, true)
			listListI32.Value = append(listListI32.Value, listI32)
			listI32.Value = listI32Value
		}
		listI64.Value = n.ListI64
	}

	var actual []byte
	actual, err = enc.Encode(&m)
	require.NoError(t, err)

	fmt.Println(expected)
	fmt.Println(actual)

	require.Equal(t, expected, actual)

}
