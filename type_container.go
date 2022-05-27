package thrift_dyn

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
)

var (
	ErrInvalidContainerItem = errors.New("invalid container type item")
)

type Container interface {
}

type Sliceable interface {
	bool | int8 | int16 | int32 | int64 |
		float64 | string | *RPCStruct | Container
}

type TypeContainerDesc struct {
	Key   thrift.TType
	Value thrift.TType
}

type TypeContainerImplementer interface {
	// Add(fs ...*TField)
	// Element() *TField

	SetSize(v int)
	GetSize() int
	// GetValue() any

	Read(ctx context.Context, p thrift.TProtocol) (err error)
	Write(ctx context.Context, p thrift.TProtocol) (err error)
}

type TypeContainer struct {
	Type  thrift.TType
	Desc  TypeContainerDesc
	Size  int
	Value reflect.Value
}

func (t *TypeContainer) Init() *TypeContainer {
	if t.Value.Kind() == reflect.Invalid {
		switch t.Type {
		case thrift.SET, thrift.LIST:
			valType := TTypeToReflectType[t.Desc.Value]
			t.Value = reflect.New(reflect.SliceOf(valType))
		case thrift.MAP:
			switch t.Desc.Key {
			case thrift.MAP, thrift.LIST, thrift.SET:
				goto InvalidType
			}
			keyType := TTypeToReflectType[t.Desc.Key]
			valType := TTypeToReflectType[t.Desc.Value]
			t.Value = reflect.MakeMapWithSize(reflect.MapOf(keyType, valType), 0)
		default:
			goto InvalidType
		}
	}
	return t
InvalidType:
	panic(fmt.Sprintf("unhandled container type"))
}

// func (t *TypeContainer) Add(fs) {
// }

func (t *TypeContainer) GetValue() any {
	return t.Value.Interface()
}
func (t *TypeContainer) SetSize(sz int) {
	t.Size = sz
}
func (t *TypeContainer) GetSize() int {
	return t.Value.Len()
}
