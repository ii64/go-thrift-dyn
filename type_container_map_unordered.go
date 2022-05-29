package thrift_dyn

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerMapUnordered[K comparable, V any] struct {
	Required bool
	Desc     TypeContainerDesc
	Size     int
	Value    map[K]V
}

func NewTypeContainerMapUnordered[K comparable, V any](desc TypeContainerDesc, required bool) *TypeContainerMapUnordered[K, V] {
	typ := &TypeContainerMapUnordered[K, V]{
		Desc:     desc,
		Required: required,
		Value:    make(map[K]V, 0),
	}
	return typ
}

func (t *TypeContainerMapUnordered[K, V]) AddKV(key K, value V) {
	t.Value[key] = value
}

func (t *TypeContainerMapUnordered[K, V]) ToMap() map[K]V {
	return t.Value
}

// ToMapPtr

func (t *TypeContainerMapUnordered[K, V]) FromMap(m map[K]V) {
	t.Value = m
}

// FromMapOrdered

func (t *TypeContainerMapUnordered[K, V]) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteMapBegin(ctx, t.Desc.Key, t.Desc.Value, size); err != nil {
		return
	}
	spec := TDataSpec{
		Required: t.Required,
		Protocol: p,
	}
	for k, v := range t.Value {
		spec.Type = t.Desc.Key
		if err = WriteData(ctx, TData[K]{
			TDataSpec: spec,
			Value:     &k,
		}); err != nil {
			return
		}

		spec.Type = t.Desc.Value
		if err = WriteData(ctx, TData[V]{
			TDataSpec: spec,
			Value:     &v,
		}); err != nil {
			return
		}
	}
	return
}

func (t *TypeContainerMapUnordered[K, V]) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	vv := map[K]V{} // new map header :/
	spec := TDataSpec{
		Required: t.Required,
		Protocol: p,
	}
	for i := 0; i < t.Size; i++ {
		var key K
		var value V
		spec.Type = t.Desc.Key
		if err = ReadData(ctx, TData[K]{
			TDataSpec: spec,
			Value:     &key,
		}); err != nil {
			return
		}

		spec.Type = t.Desc.Value
		if err = ReadData(ctx, TData[V]{
			TDataSpec: spec,
			Value:     &value,
		}); err != nil {
			return
		}
		vv[key] = value
	}
	t.Value = vv
	return
}

func (t *TypeContainerMapUnordered[K, V]) SetSize(v int) {
	t.Size = v
}
func (t *TypeContainerMapUnordered[K, V]) GetSize() int {
	return len(t.Value)
}

func NewTypeContainerMapUnorderedOfTType(desc TypeContainerDesc, required bool) (TypeContainerImplementer, error) {
	switch desc.Key {
	case thrift.BOOL:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[bool, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[bool, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[bool, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[bool, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[bool, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[bool, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[bool, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[bool, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[bool, TypeContainerImplementer](desc, required), nil
		}
	case thrift.BYTE:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[int8, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[int8, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[int8, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[int8, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[int8, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[int8, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[int8, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[int8, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[int8, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I16:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[int16, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[int16, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[int16, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[int16, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[int16, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[int16, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[int16, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[int16, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[int16, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I32:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[int32, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[int32, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[int32, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[int32, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[int32, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[int32, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[int32, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[int32, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[int32, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I64:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[int64, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[int64, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[int64, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[int64, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[int64, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[int64, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[int64, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[int64, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[int64, TypeContainerImplementer](desc, required), nil
		}
	case thrift.DOUBLE:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[float64, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[float64, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[float64, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[float64, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[float64, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[float64, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[float64, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[float64, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[float64, TypeContainerImplementer](desc, required), nil
		}
	case thrift.STRING:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMapUnordered[string, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMapUnordered[string, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMapUnordered[string, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMapUnordered[string, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMapUnordered[string, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMapUnordered[string, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMapUnordered[string, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMapUnordered[string, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMapUnordered[string, TypeContainerImplementer](desc, required), nil
		}
	}
	return nil, fmt.Errorf("unhandled type key:{%s} value:{%s}", desc.Key, desc.Value)
}
