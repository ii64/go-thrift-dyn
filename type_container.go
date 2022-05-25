package thrift_dyn

import (
	"context"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
)

var (
	ErrInvalidContainerItem = errors.New("invalid container type item")
)

type TypeContainerDesc struct {
	Key   thrift.TType
	Value thrift.TType
}

type TypeContainer struct {
	Type  thrift.TType
	Desc  TypeContainerDesc
	Size  int
	Value []*TField
}

type TypeContainerImplementer interface {
	Add(fs ...*TField) error
	Element() *TField

	SetSize(v int)
	GetSize() int
	GetValue() any

	Read(ctx context.Context, p thrift.TProtocol) (err error)
	Write(ctx context.Context, p thrift.TProtocol) (err error)
}

func (t *TypeContainer) Add(fs ...*TField) error {
	t.Value = append(t.Value, fs...)
	return nil
}

func (t *TypeContainer) GetValue() any {
	return t.Value
}
func (t *TypeContainer) SetSize(sz int) {
	t.Size = sz
}
func (t *TypeContainer) GetSize() int {
	return len(t.Value)
}
