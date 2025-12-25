package documentstore

import (
	"errors"
	"fmt"
	"reflect"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value any
}

type Document struct {
	Fields map[string]DocumentField
}

func MarshalDocument(input any) (*Document, error) {
	if input == nil {
		return nil, errors.New("MarshalDocument: input is nil")
	}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, errors.New("MarshalDocument: input pointer is nil")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("MarshalDocument: expected struct or *struct, got %s", v.Kind())
	}

	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)

		// Skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		fieldName := sf.Name

		fieldVal := v.Field(i)
		df, err := marshalValue(fieldVal)
		if err != nil {
			return nil, fmt.Errorf("MarshalDocument: field %q: %w", fieldName, err)
		}

		doc.Fields[fieldName] = df
	}

	return doc, nil
}

func marshalValue(v reflect.Value) (DocumentField, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return DocumentField{
				Type:  DocumentFieldTypeObject,
				Value: (*Document)(nil),
			}, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return DocumentField{
			Type:  DocumentFieldTypeString,
			Value: v.String(),
		}, nil

	case reflect.Bool:
		return DocumentField{
			Type:  DocumentFieldTypeBool,
			Value: v.Bool(),
		}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return DocumentField{
			Type:  DocumentFieldTypeNumber,
			Value: v.Interface(),
		}, nil

	case reflect.Slice, reflect.Array:
		n := v.Len()
		items := make([]DocumentField, 0, n)
		for i := range n {
			elem := v.Index(i)
			df, err := marshalValue(elem)
			if err != nil {
				return DocumentField{}, fmt.Errorf("array element %d: %w", i, err)
			}
			items = append(items, df)
		}
		return DocumentField{
			Type:  DocumentFieldTypeArray,
			Value: items,
		}, nil

	case reflect.Struct:
		nested, err := MarshalDocument(v.Interface())
		if err != nil {
			return DocumentField{}, err
		}
		return DocumentField{
			Type:  DocumentFieldTypeObject,
			Value: nested,
		}, nil

	default:
		return DocumentField{}, fmt.Errorf("unsupported kind %s", v.Kind())
	}
}

func UnmarshalDocument(doc *Document, output any) error {
	if doc == nil {
		return errors.New("UnmarshalDocument: doc is nil")
	}
	if output == nil {
		return errors.New("UnmarshalDocument: output is nil")
	}

	v := reflect.ValueOf(output)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("UnmarshalDocument: output must be a non-nil pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("UnmarshalDocument: expected pointer to struct, got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		// Skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		fieldName := sf.Name
		df, ok := doc.Fields[fieldName]
		if !ok {
			// Missing from document: leave zero value
			continue
		}

		destField := v.Field(i)
		if !destField.CanSet() {
			continue
		}

		if err := unmarshalValue(df, destField); err != nil {
			return fmt.Errorf("field %q: %w", fieldName, err)
		}
	}

	return nil
}

func unmarshalValue(df DocumentField, dest reflect.Value) error {
	if dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		return unmarshalValue(df, dest.Elem())
	}

	switch dest.Kind() {
	case reflect.String:
		if df.Type != DocumentFieldTypeString {
			return fmt.Errorf("expected string, got %s", df.Type)
		}
		s, ok := df.Value.(string)
		if !ok {
			return fmt.Errorf("stored value is not string, got %T", df.Value)
		}
		dest.SetString(s)
		return nil

	case reflect.Bool:
		if df.Type != DocumentFieldTypeBool {
			return fmt.Errorf("expected bool, got %s", df.Type)
		}
		b, ok := df.Value.(bool)
		if !ok {
			return fmt.Errorf("stored value is not bool, got %T", df.Value)
		}
		dest.SetBool(b)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		if df.Type != DocumentFieldTypeNumber {
			return fmt.Errorf("expected number, got %s", df.Type)
		}
		fv := reflect.ValueOf(df.Value)
		if !fv.IsValid() {
			return fmt.Errorf("number value is invalid")
		}
		if !fv.Type().ConvertibleTo(dest.Type()) {
			return fmt.Errorf("cannot convert %T to %s", df.Value, dest.Type())
		}
		dest.Set(fv.Convert(dest.Type()))
		return nil

	case reflect.Slice:
		if df.Type != DocumentFieldTypeArray {
			return fmt.Errorf("expected array, got %s", df.Type)
		}
		items, ok := df.Value.([]DocumentField)
		if !ok {
			return fmt.Errorf("stored value is not []DocumentField, got %T", df.Value)
		}

		slice := reflect.MakeSlice(dest.Type(), len(items), len(items))
		for i, item := range items {
			if err := unmarshalValue(item, slice.Index(i)); err != nil {
				return fmt.Errorf("array element %d: %w", i, err)
			}
		}
		dest.Set(slice)
		return nil

	case reflect.Array:
		if df.Type != DocumentFieldTypeArray {
			return fmt.Errorf("expected array, got %s", df.Type)
		}
		items, ok := df.Value.([]DocumentField)
		if !ok {
			return fmt.Errorf("stored value is not []DocumentField, got %T", df.Value)
		}
		if len(items) != dest.Len() {
			return fmt.Errorf("array length mismatch: have %d, need %d", len(items), dest.Len())
		}
		for i, item := range items {
			if err := unmarshalValue(item, dest.Index(i)); err != nil {
				return fmt.Errorf("array element %d: %w", i, err)
			}
		}
		return nil

	case reflect.Struct:
		if df.Type != DocumentFieldTypeObject {
			return fmt.Errorf("expected object, got %s", df.Type)
		}
		if df.Value == nil {
			return nil
		}
		nestedDoc, ok := df.Value.(*Document)
		if !ok {
			return fmt.Errorf("stored value is not *Document, got %T", df.Value)
		}
		// dest is a struct; we need its address to pass to UnmarshalDocument
		return UnmarshalDocument(nestedDoc, dest.Addr().Interface())

	default:
		return fmt.Errorf("unsupported destination kind %s", dest.Kind())
	}
}
