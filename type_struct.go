package go_thriftproxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type RPCStruct struct {
	Name   string
	Fields []*TField
}

func (s *RPCStruct) AddField(fs ...*TField) *RPCStruct {
	s.Fields = append(s.Fields, fs...)
	return s
}

// Write writes fields to the wire
func (s *RPCStruct) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	var fieldId TFieldID
	if err = p.WriteStructBegin(ctx, s.Name); err != nil {
		goto WriteStructBeginError
	}

	for _, field := range s.Fields {
		fieldId = field.GetID()
		if err = field.Write(ctx, p); err != nil {
			goto WriteFieldError
		}
	}

	if err = p.WriteFieldStop(ctx); err != nil {
		goto WriteFieldStopError
	}
	if err = p.WriteStructEnd(ctx); err != nil {
		goto WriteStructEndError
	}
	return nil
WriteStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", s.Name), err)
WriteFieldError:
	return thrift.PrependError(fmt.Sprintf("%T write field %d error: ", p, fieldId), err)
WriteFieldStopError:
	return thrift.PrependError(fmt.Sprintf("%T write field stop error: ", s.Name), err)
WriteStructEndError:
	return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", s.Name), err)
}

// Read reads fields from wire
func (s *RPCStruct) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	var (
		fieldName   string
		fieldTypeId thrift.TType
		fieldId     TFieldID
		vv          = s.Fields[:0]
	)

	if _, err = p.ReadStructBegin(ctx); err != nil {
		goto ReadStructBeginError
	}

	for {
		fieldName, fieldTypeId, fieldId, err = p.ReadFieldBegin(ctx)
		if err != nil {
			goto ReadFieldBeginError
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		var field = NewTField(fieldId, fieldTypeId, fieldName, false)
		if err = field.Read(ctx, p); errors.Is(err, ErrSkipField) {
			if err = p.Skip(ctx, fieldTypeId); err != nil {
				goto SkipFieldError
			}
		} else if err != nil {
			goto ReadFieldError
		}

		if err = p.ReadFieldEnd(ctx); err != nil {
			goto ReadFieldEndError
		}

		vv = append(vv, field)
	}
	s.Fields = vv

	if err = p.ReadStructEnd(ctx); err != nil {
		goto ReadStructEndError
	}

	return nil
ReadStructBeginError:
	return thrift.PrependError(fmt.Sprintf("%s read struct begin error: ", s.Name), err)
ReadFieldBeginError:
	return thrift.PrependError(fmt.Sprintf("%s read field %d begin error: ", s.Name, fieldId), err)
ReadFieldError:
	return thrift.PrependError(fmt.Sprintf("%s read field %d '%s' (%d) error: ", s.Name, fieldId, fieldName, fieldTypeId), err)
SkipFieldError:
	return thrift.PrependError(fmt.Sprintf("%s field %d skip type %d error: ", s.Name, fieldId, fieldTypeId), err)
ReadFieldEndError:
	return thrift.PrependError(fmt.Sprintf("%s read field end error", s.Name), err)
ReadStructEndError:
	return thrift.PrependError(fmt.Sprintf("%s read struct end error: ", s.Name), err)
}
