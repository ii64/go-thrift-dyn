package thrift_dyn

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerSet[T Sliceable] struct {
	Required bool
	Desc     TypeContainerDesc
	Size     int
	Value    []T
}

func NewTypeContainerSet[T Sliceable](desc TypeContainerDesc, required bool) *TypeContainerSet[T] {
	typ := &TypeContainerSet[T]{
		Required: required,
		Desc:     desc,
	}
	return typ
}

func (t *TypeContainerSet[T]) Add(vs ...T) {
	t.Value = append(t.Value, vs...)
}

func (t *TypeContainerSet[T]) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteSetBegin(ctx, t.Desc.Value, size); err != nil {
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
	if err = p.WriteSetEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerSet[T]) Read(ctx context.Context, p thrift.TProtocol) (err error) {
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
		if err = ReadData[T](ctx, data); err != nil {
			return
		}
		vv = append(vv, value)
	}
	t.Value = vv
	return
}

func (t *TypeContainerSet[T]) SetSize(v int) {
	t.Size = v
}
func (t *TypeContainerSet[T]) GetSize() int {
	return len(t.Value)
}

func NewTypeContainerSetOfTType(desc TypeContainerDesc, required bool) (TypeContainerImplementer, error) {
	switch desc.Value {
	case thrift.BOOL:
		return NewTypeContainerSet[bool](desc, required), nil
	case thrift.BYTE:
		return NewTypeContainerSet[int8](desc, required), nil
	case thrift.I16:
		return NewTypeContainerSet[int16](desc, required), nil
	case thrift.I32:
		return NewTypeContainerSet[int32](desc, required), nil
	case thrift.I64:
		return NewTypeContainerSet[int64](desc, required), nil
	case thrift.DOUBLE:
		return NewTypeContainerSet[float64](desc, required), nil
	case thrift.STRING:
		return NewTypeContainerSet[string](desc, required), nil
	case thrift.STRUCT:
		return NewTypeContainerSet[thrift.TStruct](desc, required), nil
	case thrift.MAP, thrift.LIST, thrift.SET:
		return NewTypeContainerSet[TypeContainerImplementer](desc, required), nil
	}
	return nil, fmt.Errorf("unhandled type %T", desc.Value)
}
