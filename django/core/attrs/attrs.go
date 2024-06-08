package attrs

import (
	"github.com/Nigel2392/django/forms/fields"
)

type Definer interface {
	FieldDefs() Definitions
}

type Definitions interface {
	Set(name string, value interface{}) error
	Get(name string) interface{}
	Field(name string) (f Field, ok bool)
	ForceSet(name string, value interface{}) error
	Primary() Field
	Fields() []Field
}

type Field interface {
	Labeler
	Helper
	Stringer
	Namer
	Instance() Definer
	IsPrimary() bool
	AllowNull() bool
	AllowBlank() bool
	AllowEdit() bool
	GetValue() interface{}
	GetDefault() interface{}
	SetValue(v interface{}, force bool) error
	FormField() fields.Field
	Validate() error
}

type Namer interface {
	Name() string
}

type Stringer interface {
	ToString() string
}

type Labeler interface {
	Label() string
}

type Helper interface {
	HelpText() string
}
