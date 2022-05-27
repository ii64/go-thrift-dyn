package thrift_dyn

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"golang.org/x/exp/slices"
)

type TypeContainerMapItem[K comparable, V any] struct {
	Key   K
	Value V
}

type TypeContainerMap[K comparable, V any] struct {
	Type     thrift.TType
	Required bool
	Desc     TypeContainerDesc
	Size     int
	Value    []TypeContainerMapItem[K, V]
}

func NewTypeContainerMap[K comparable, V any](desc TypeContainerDesc, required bool) *TypeContainerMap[K, V] {
	typ := &TypeContainerMap[K, V]{
		Type:     thrift.MAP,
		Required: required,
		Desc:     desc,
	}
	return typ
}

func (t *TypeContainerMap[K, V]) Add(vs ...TypeContainerMapItem[K, V]) {
	t.Value = append(t.Value, vs...)
}
func (t *TypeContainerMap[K, V]) AddKV(key K, value V) {
	t.Value = append(t.Value, TypeContainerMapItem[K, V]{
		Key:   key,
		Value: value,
	})
}

func (t *TypeContainerMap[K, V]) ToMap() map[K]V {
	var ret = map[K]V{}
	for _, item := range t.Value {
		ret[item.Key] = item.Value
	}
	return ret
}

func (t *TypeContainerMap[K, V]) ToMapPtr() map[K]*V {
	var ret = map[K]*V{}
	for i := range t.Value {
		key := t.Value[i].Key
		value := &t.Value[i].Value
		ret[key] = value
	}
	return ret
}

func (t *TypeContainerMap[K, V]) FromMap(m map[K]V) {
	for k, v := range m {
		t.AddKV(k, v)
	}
}
func (t *TypeContainerMap[K, V]) FromMapOrdered(m map[K]V) {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	t.mapKeySorter(keys)
	for _, k := range keys {
		v := m[k]
		t.AddKV(k, v)
	}
}

func (t *TypeContainerMap[K, V]) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteMapBegin(ctx, t.Desc.Key, t.Desc.Value, size); err != nil {
		return
	}
	for i := 0; i < size; i++ {
		value := t.Value[i]
		if err = WriteData[K](ctx, TData[K]{
			TDataSpec: TDataSpec{
				Type:     t.Desc.Key,
				Required: t.Required,
				Protocol: p,
			},
			Value: &value.Key,
		}); err != nil {
			return
		}
		if err = WriteData[V](ctx, TData[V]{
			TDataSpec: TDataSpec{
				Type:     t.Desc.Value,
				Required: t.Required,
				Protocol: p,
			},
			Value: &value.Value,
		}); err != nil {
			return
		}
	}
	if err = p.WriteMapEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerMap[K, V]) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	vv := t.Value[:0]
	for i := 0; i < t.Size; i++ {
		var value TypeContainerMapItem[K, V]
		spec := TDataSpec{
			Required: t.Required,
			Protocol: p,
		}

		spec.Type = t.Desc.Key
		if err = ReadData(ctx, TData[K]{
			TDataSpec: spec,
			Value:     &value.Key,
		}); err != nil {
			return
		}

		spec.Type = t.Desc.Value
		if err = ReadData(ctx, TData[V]{
			TDataSpec: spec,
			Value:     &value.Value,
		}); err != nil {
			return
		}

		vv = append(vv, value)
	}
	t.Value = vv
	return
}

func (t *TypeContainerMap[K, V]) SetSize(v int) {
	t.Size = v
}
func (t *TypeContainerMap[K, V]) GetSize() int {
	return len(t.Value)
}

func (t *TypeContainerMap[K, V]) mapKeySorter(keys []K) {
	switch v := (any)(keys).(type) {
	case []bool:
		slices.SortFunc(v, func(a, b bool) bool {
			if a {
				return true
			}
			return false
		})
	case []int8:
		slices.Sort(v)
	case []int16:
		slices.Sort(v)
	case []int32:
		slices.Sort(v)
	case []int64:
		slices.Sort(v)
	case []float64:
		slices.Sort(v)
	case []string:
		slices.Sort(v)
	}
}

func NewTypeContainerMapOfTType(desc TypeContainerDesc, required bool) (TypeContainerImplementer, error) {
	switch desc.Key {
	case thrift.BOOL:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[bool, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[bool, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[bool, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[bool, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[bool, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[bool, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[bool, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[bool, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[bool, TypeContainerImplementer](desc, required), nil
		}
	case thrift.BYTE:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[int8, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[int8, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[int8, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[int8, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[int8, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[int8, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[int8, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[int8, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[int8, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I16:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[int16, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[int16, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[int16, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[int16, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[int16, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[int16, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[int16, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[int16, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[int16, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I32:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[int32, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[int32, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[int32, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[int32, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[int32, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[int32, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[int32, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[int32, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[int32, TypeContainerImplementer](desc, required), nil
		}
	case thrift.I64:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[int64, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[int64, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[int64, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[int64, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[int64, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[int64, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[int64, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[int64, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[int64, TypeContainerImplementer](desc, required), nil
		}
	case thrift.DOUBLE:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[float64, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[float64, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[float64, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[float64, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[float64, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[float64, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[float64, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[float64, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[float64, TypeContainerImplementer](desc, required), nil
		}
	case thrift.STRING:
		switch desc.Value {
		case thrift.BOOL:
			return NewTypeContainerMap[string, bool](desc, required), nil
		case thrift.BYTE:
			return NewTypeContainerMap[string, int8](desc, required), nil
		case thrift.I16:
			return NewTypeContainerMap[string, int16](desc, required), nil
		case thrift.I32:
			return NewTypeContainerMap[string, int32](desc, required), nil
		case thrift.I64:
			return NewTypeContainerMap[string, int64](desc, required), nil
		case thrift.DOUBLE:
			return NewTypeContainerMap[string, float64](desc, required), nil
		case thrift.STRING:
			return NewTypeContainerMap[string, string](desc, required), nil
		case thrift.STRUCT:
			return NewTypeContainerMap[string, thrift.TStruct](desc, required), nil
		case thrift.MAP, thrift.LIST, thrift.SET:
			return NewTypeContainerMap[string, TypeContainerImplementer](desc, required), nil
		}
	}
	return nil, fmt.Errorf("unhandled type key:{%s} value:{%s}", desc.Key, desc.Value)
}
