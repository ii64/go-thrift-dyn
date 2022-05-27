package thrift_dyn

import (
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
	"unsafe"
)

const _BT_SIZE = 16
const (
	_BT_F_RTTYPE = 1 << 0
	_BT_F_RTBRW  = 1 << 1
)

var (
	_BT_INIT           int
	TTypeToReflectType [_BT_SIZE]reflect.Type
)

func init() {
	if _BT_INIT&_BT_F_RTTYPE == 0 {
		TTypeToReflectType[thrift.STOP] = nil
		TTypeToReflectType[thrift.VOID] = nil
		TTypeToReflectType[thrift.BOOL] = reflect.TypeOf(true)
		TTypeToReflectType[thrift.BYTE] = reflect.TypeOf(int8(0)) // BYTE, I08
		TTypeToReflectType[thrift.DOUBLE] = reflect.TypeOf(float64(0))
		TTypeToReflectType[thrift.I16] = reflect.TypeOf(int16(0))
		TTypeToReflectType[thrift.I32] = reflect.TypeOf(int32(0))
		TTypeToReflectType[thrift.I64] = reflect.TypeOf(int64(0))
		TTypeToReflectType[thrift.STRING] = reflect.TypeOf("")
		TTypeToReflectType[thrift.STRUCT] = reflect.TypeOf(&RPCStruct{})
		// Container type
		// TTypeToReflectType[thrift.SET] = reflect.TypeOf(TypeContainerImplementer(&TypeContainerSet{}))
		// TTypeToReflectType[thrift.LIST] = reflect.TypeOf(TypeContainerImplementer(&TypeContainerList{}))
		// TTypeToReflectType[thrift.MAP] = reflect.TypeOf(TypeContainerImplementer(&TypeContainerMap{}))

		_BT_INIT |= _BT_F_RTTYPE
	}
}

func typOff(typ_, ktyp, vtyp thrift.TType) {

}

func mapassign[K comparable, V any](dst any, key K, value V) {
	word := (*[2]unsafe.Pointer)(unsafe.Pointer(&dst))[1]
	tmp := *(*map[K]V)(unsafe.Pointer(&word))
	tmp[key] = value
}

func mapaccess[K comparable, V any](src any, key K, value *V) {
	word := (*[2]unsafe.Pointer)(unsafe.Pointer(&src))[1]
	tmp := *(*map[K]V)(unsafe.Pointer(&word))
	*value = tmp[key]
}

func String2bs(s string) (r []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sb := (*reflect.SliceHeader)(unsafe.Pointer(&r))
	sb.Data = sh.Data
	sb.Len, sb.Cap = sh.Len, sh.Len
	return
}
