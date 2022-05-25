package go_thriftproxy

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
)

type TypeContainerKV struct {
	Key   *TField
	Value *TField
}

func (kv *TypeContainerKV) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	if thunk := kv.Key.WriteDataThunk(); thunk != nil {
		if err = thunk(ctx, p); err != nil {
			return
		}
	}
	if thunk := kv.Value.WriteDataThunk(); thunk != nil {
		if err = thunk(ctx, p); err != nil {
			return
		}
	}
	return
}

func (kv *TypeContainerKV) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	if err = kv.Key.Read(ctx, p); err != nil {
		return
	}
	if err = kv.Value.Read(ctx, p); err != nil {
		return
	}
	return
}

type TypeContainerMap struct {
	TypeContainer
}

func NewTypeContainerMap(type_ thrift.TType, desc TypeContainerDesc, required bool) TypeContainerImplementer {
	return &TypeContainerMap{
		TypeContainer: TypeContainer{
			Type: type_, Desc: desc,
		},
	}
}

func (t *TypeContainerMap) newElement(type_ thrift.TType) *TField {
	f := NewTField(0, type_, "map_item", true)
	f.SetValue(&TypeContainerKV{
		Key:   NewTField(0, t.Desc.Key, "map_key", true),
		Value: NewTField(0, t.Desc.Value, "map_value", true),
	})
	return f
}
func (t *TypeContainerMap) Element() *TField {
	return t.newElement(0)
}

func (t *TypeContainerMap) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	size := t.GetSize()
	if err = p.WriteMapBegin(ctx, t.Desc.Key, t.Desc.Value, size); err != nil {
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
	if err = p.WriteMapEnd(ctx); err != nil {
		return
	}
	return
}

func (t *TypeContainerMap) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	var vv = t.Value[:0]
	for i := 0; i < t.Size; i++ {
		item := t.Element()
		if err = item.Read(ctx, p); err != nil {
			return
		}
		vv = append(vv, item)
	}
	t.Value = vv
	return
}
