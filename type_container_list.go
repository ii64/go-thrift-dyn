package go_thriftproxy

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerList struct {
	TypeContainer
}

func NewTypeContainerList(type_ thrift.TType, desc TypeContainerDesc, required bool) TypeContainerImplementer {
	return &TypeContainerList{
		TypeContainer: TypeContainer{
			Type: type_, Desc: desc,
		},
	}
}

func (t *TypeContainerList) newElement(type_ thrift.TType) *TField {
	return NewTField(0, type_, "list_value", true)
}
func (t *TypeContainerList) Element() *TField {
	return t.newElement(t.Desc.Value)
}

func (t *TypeContainerList) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteListBegin(ctx, t.Desc.Value, size); err != nil {
		return
	}
	for i := 0; i < size; i++ {
		value := t.Value[i]
		if thunk := value.WriteDataThunk(); thunk != nil {
			if err = thunk(ctx, p); err != nil {
				return
			}
		}
	}
	if err = p.WriteListEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerList) Read(ctx context.Context, p thrift.TProtocol) (err error) {
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
