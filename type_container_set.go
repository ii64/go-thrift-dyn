package go_thriftproxy

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerSet struct {
	TypeContainer
}

func NewTypeContainerSet(type_ thrift.TType, desc TypeContainerDesc, required bool) TypeContainerImplementer {
	return &TypeContainerSet{
		TypeContainer: TypeContainer{
			Type: type_, Desc: desc,
		},
	}
}

func (t *TypeContainerSet) newElement(type_ thrift.TType) *TField {
	return NewTField(0, type_, "set_value", true)
}
func (t *TypeContainerSet) Element() *TField {
	return t.newElement(t.Desc.Value)
}

func (t *TypeContainerSet) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteSetBegin(ctx, t.Desc.Value, size); err != nil {
		return
	}
	for i := 0; i < size; i++ {
		value := t.Value[i]
		if err = value.WriteData(ctx, p); err != nil {
			return
		}
	}
	if err = p.WriteSetEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerSet) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	var vv = t.Value[:0]
	for i := 0; i < t.Size; i++ {
		el := t.Element()
		if err = el.Read(ctx, p); err != nil {
			return
		}
		vv = append(vv, el)
	}
	t.Value = vv
	return
}
