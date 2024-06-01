package attrs

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"reflect"
	"time"

	"github.com/Nigel2392/django/core/assert"
	"github.com/Nigel2392/django/forms/fields"
	"github.com/Nigel2392/django/forms/widgets"
)

type FieldDef struct {
	Blank          bool
	Null           bool
	Editable       bool
	instance_t_ptr reflect.Type
	instance_v_ptr reflect.Value
	instance_t     reflect.Type
	instance_v     reflect.Value
	field_t        reflect.StructField
	field_v        reflect.Value
}

func NewField[T any](instance *T, name string, null, blank, editable bool) *FieldDef {
	var (
		instance_t_ptr = reflect.TypeOf(instance)
		instance_v_ptr = reflect.ValueOf(instance)
		instance_t     = instance_t_ptr.Elem()
		instance_v     = instance_v_ptr.Elem()
		field_t        reflect.StructField
		field_v        reflect.Value
		ok             bool
	)

	field_t, ok = instance_t.FieldByName(name)
	assert.True(ok, "field %q not found in %T", name, instance)

	field_v = instance_v.FieldByIndex(field_t.Index)
	assert.True(field_v.IsValid(), "field %q not found in %T", name, instance)

	return &FieldDef{
		Null:           null,
		Blank:          blank,
		Editable:       editable,
		instance_t_ptr: instance_t_ptr,
		instance_v_ptr: instance_v_ptr,
		instance_t:     instance_t,
		instance_v:     instance_v,
		field_t:        field_t,
		field_v:        field_v,
	}
}

func (f *FieldDef) Label() string {
	if labeler, ok := f.field_v.Interface().(Labeler); ok {
		return labeler.Label()
	}
	return ""
}

func (f *FieldDef) HelpText() string {
	if helpTexter, ok := f.field_v.Interface().(Helper); ok {
		return helpTexter.HelpText()
	}
	return ""
}

func (f *FieldDef) Name() string {
	return f.field_t.Name
}

func (f *FieldDef) AllowNull() bool {
	return f.Null
}

func (f *FieldDef) AllowBlank() bool {
	return f.Blank
}

func (f *FieldDef) AllowEdit() bool {
	return f.Editable
}

func (f *FieldDef) Validate() error {
	return nil
}

func (f *FieldDef) GetValue() interface{} {
	return f.field_v.Interface()
}

func (f *FieldDef) GetDefault() interface{} {

	var funcName = fmt.Sprintf("GetDefault%s", f.Name())
	if method, ok := f.instance_t.MethodByName(funcName); ok {
		return method.Func.Call([]reflect.Value{f.instance_v_ptr})[0].Interface()
	}

	if !f.field_v.IsValid() {
		return reflect.Zero(f.field_t.Type).Interface()
	}

	return f.field_v.Interface()
}

func (f *FieldDef) FormField() fields.Field {
	var opts = make([]func(fields.Field), 0)
	if f.Label() != "" {
		opts = append(opts, fields.Label(f.Label()))
	}
	opts = append(opts,
		fields.Name(f.Name()),
	)

	switch reflect.New(f.field_t.Type).Elem().Interface().(type) {
	case time.Time:
		return fields.DateField(widgets.DateWidgetTypeDateTime, opts...)
	case json.RawMessage:
		return fields.JSONField[map[string]interface{}](opts...)
	case mail.Address, *mail.Address:
		return fields.EmailField(opts...)
	}

	switch f.field_t.Type.Kind() {
	case reflect.String:
		return fields.CharField(opts...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fields.NumberField[int](opts...)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fields.NumberField[uint](opts...)
	case reflect.Float32, reflect.Float64:
		return fields.NumberField[float64](opts...)
	default:
		return fields.CharField(opts...)
	}
}

func (f *FieldDef) SetValue(v interface{}, force bool) error {
	var r_v = reflect.ValueOf(v)

	if err := assert.True(
		r_v.IsValid() || f.AllowNull(),
		"field %q (%q) is not valid", f.field_t.Name, f.field_t.Type,
	); err != nil {
		return err
	}

	if !r_v.IsValid() && f.AllowNull() {
		f.field_v.Set(reflect.Zero(f.field_t.Type))
		return nil
	}

	if err := assert.True(
		r_v.Type() == f.field_t.Type || r_v.CanConvert(f.field_t.Type),
		"field %q (%q) is not convertible to %q",
		f.field_t.Name, r_v.Type(), f.field_t.Type,
	); err != nil {
		return err
	}

	if r_v.Type() != f.field_t.Type {
		r_v = r_v.Convert(f.field_t.Type)
	}

	if err := assert.True(
		f.field_v.CanSet() && (f.Editable || force),
		"field %q is not editable", f.field_t.Name,
	); err != nil {
		return err
	}

	if r_v.IsZero() && !f.AllowBlank() {
		switch r_v.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128:
		default:
			return assert.Fail(
				fmt.Sprintf("field %q is not blank", f.field_t.Name),
			)
		}
	}

	f.field_v.Set(r_v)
	return nil
}
