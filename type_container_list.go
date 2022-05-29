package thrift_dyn

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerList[T Sliceable] struct {
	Required bool
	Desc     TypeContainerDesc
	Size     int
	Value    []T
}

func NewTypeContainerList[T Sliceable](desc TypeContainerDesc, required bool) *TypeContainerList[T] {
	typ := &TypeContainerList[T]{
		Required: required,
		Desc:     desc,
	}
	return typ
}

func (t *TypeContainerList[T]) Add(vs ...T) {
	t.Value = append(t.Value, vs...)
}

func (t *TypeContainerList[T]) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteListBegin(ctx, t.Desc.Value, size); err != nil {
		return
	}
	for i := 0; i < size; i++ {
		if err = WriteData[T](ctx, TData[T]{
			TDataSpec: TDataSpec{
				Type:     t.Desc.Value,
				Required: t.Required,
				Protocol: p,
			},
			Value: &t.Value[i],
		}); err != nil {
			return
		}
	}
	if err = p.WriteListEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerList[T]) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	vv := t.Value[:0]
	for i := 0; i < t.Size; i++ {
		var value T
		data := TData[T]{
			TDataSpec: TDataSpec{
				Type:     t.Desc.Value,
				Required: t.Required,
				Protocol: p,
			},
			Value: &value,
		}
		if err = ReadData(ctx, data); err != nil {
			return
		}
		vv = append(vv, value)
	}
	t.Value = vv
	return
}

func (t *TypeContainerList[T]) SetSize(v int) {
	t.Size = v
}
func (t *TypeContainerList[T]) GetSize() int {
	return len(t.Value)
}

func NewTypeContainerListOfTType(desc TypeContainerDesc, required bool) (TypeContainerImplementer, error) {
	switch desc.Value {
	case thrift.BOOL:
		return NewTypeContainerList[bool](desc, required), nil
	case thrift.BYTE:
		return NewTypeContainerList[int8](desc, required), nil
	case thrift.I16:
		return NewTypeContainerList[int16](desc, required), nil
	case thrift.I32:
		return NewTypeContainerList[int32](desc, required), nil
	case thrift.I64:
		return NewTypeContainerList[int64](desc, required), nil
	case thrift.DOUBLE:
		return NewTypeContainerList[float64](desc, required), nil
	case thrift.STRING:
		return NewTypeContainerList[string](desc, required), nil
	case thrift.STRUCT:
		return NewTypeContainerList[thrift.TStruct](desc, required), nil
	case thrift.MAP, thrift.LIST, thrift.SET:
		return NewTypeContainerList[TypeContainerImplementer](desc, required), nil
	}
	return nil, fmt.Errorf("unhandled type %T", desc.Value)
}
