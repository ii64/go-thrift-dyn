package thrift_dyn

import (
	"context"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
)

var (
	ErrSkipField = errors.New("skip field")
)

type TFieldID = int16

type TField struct {
	ID       TFieldID     // i16
	Type     thrift.TType // i8
	Required bool         // i8
	Name     string

	Value any
}

type TFieldImplementerx interface {
	GetID() TFieldID
	GetValue() any

	WriteDataThunk(ctx context.Context, p thrift.TProtocol) (doFunc func() error)
	Write(ctx context.Context, p thrift.TProtocol) (err error)
	Read(ctx context.Context, p thrift.TProtocol) (err error)
}

// NewTField create new TField.
func NewTField(id TFieldID, type_ thrift.TType, name string, required bool) *TField {
	return &TField{ID: id, Name: name, Type: type_, Required: required}
}

// WriteData write data to wire, without writing struct field.
func (f *TField) WriteData(ctx context.Context, p thrift.TProtocol) (err error) {
	return WriteDataGeneric(ctx, TDataSpec{
		Type:     f.Type,
		Required: f.Required,
		Protocol: p,
	}, f.Value)
}

// Write write struct field with its value.
func (f *TField) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	if f.Value != nil || f.Required {
		if err = p.WriteFieldBegin(ctx, f.Name, f.Type, f.ID); err != nil {
			return
		}
		if err = f.WriteData(ctx, p); err != nil {
			return
		}
		if err = p.WriteFieldEnd(ctx); err != nil {
			return
		}
	}
	return
}

// Read read struct field with its value.
func (f *TField) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	spec := TDataSpec{
		Type:     f.Type,
		Required: f.Required,
		Protocol: p,
	}
	if err = ReadDataGeneric(ctx, spec, &f.Value); err != nil {
		// ErrSkipField
		return
	}
	return
}

// GetID get field ID.
func (f TField) GetID() TFieldID {
	return f.ID
}

// GetValue get field value.
func (f *TField) GetValue() any {
	return f.Value
}

// SetValue set value of field, or value of KV container (map)
func (f *TField) SetValue(value any) *TField {
	f.Value = value
	return f
}
