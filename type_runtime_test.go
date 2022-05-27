package thrift_dyn

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestBasicTTypeToReflectTypeSlice(t *testing.T) {
	for ttype, btype := range TTypeToReflectType {
		_ = ttype
		if btype == nil {
			continue
		}
		allocSlice := reflect.New(reflect.SliceOf(btype))
		switch ttype {
		case thrift.BOOL:
		case thrift.BYTE:
		case thrift.DOUBLE:
		case thrift.I16:
		case thrift.I32:
		case thrift.I64:
		case thrift.STRING:
		case thrift.STRUCT:
		}

		pt := (*reflect.SliceHeader)(allocSlice.UnsafePointer())
		fmt.Printf("%+#v %+#v\n", allocSlice, pt)
	}

}

func TestString2bs(t *testing.T) {
	s := "hello world 世界"
	bs := String2bs(s)
	require.Equal(t, bs, []byte(s))
}
