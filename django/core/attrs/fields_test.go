package attrs_test

import (
	"reflect"
	"testing"

	"github.com/Nigel2392/django/core/attrs"
	"github.com/Nigel2392/django/forms/fields"
	"github.com/Nigel2392/django/forms/widgets"
	"github.com/Nigel2392/goldcrest"
)

type TestModelFields struct {
	ID      int
	Name    string
	Objects []int64
}

type customTestWidget struct {
	*widgets.BaseWidget
}

func (f *TestModelFields) FieldDefs() attrs.Definitions {
	return attrs.Define(f,
		attrs.NewField(f, "ID", false, false, true),
		attrs.NewField(f, "Name", false, false, true),
		attrs.NewField(f, "Objects", false, false, false),
	)
}

func TestModelFieldsGet(t *testing.T) {
	var m = &TestModelDefinitions{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, true)
	)

	if m.ID != defID.GetValue().(int) {
		t.Errorf("expected %d, got %d", m.ID, defID.GetValue())
	}

	if m.Name != defName.GetValue().(string) {
		t.Errorf("expected %q, got %q", m.Name, defName.GetValue())
	}

	if len(m.Objects) != len(defObjects.GetValue().([]int64)) {
		t.Errorf("expected %d, got %d", len(m.Objects), len(defObjects.GetValue().([]int64)))
	}
}

func TestModelFieldFieldsSet(t *testing.T) {
	var m = &TestModelDefinitions{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, true)
	)

	defID.SetValue(2, false)
	defName.SetValue("new name", false)
	defObjects.SetValue([]int64{4, 5, 6}, false)

	if m.ID != 2 {
		t.Errorf("expected %d, got %d", 2, m.ID)
	}

	if m.Name != "new name" {
		t.Errorf("expected %q, got %q", "new name", m.Name)
	}

	if len(m.Objects) != 3 {
		t.Errorf("expected %d, got %d", 3, len(m.Objects))
	}

	if m.Objects[0] != 4 {
		t.Errorf("expected %d, got %d", 4, m.Objects[0])
	}
}

func TestModelFieldFieldsSetReadOnly(t *testing.T) {
	var m = &TestModelFields{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, false)
	)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got nil")
		}

		if m.Objects[0] != 1 {
			t.Errorf("expected %d, got %d", 1, m.Objects[0])
		}

		if m.ID != 2 {
			t.Errorf("expected %d, got %d", 1, m.ID)
		}

		if m.Name != "new name" {
			t.Errorf("expected %q, got %q", "name", m.Name)
		}
	}()

	defID.SetValue(2, false)
	defName.SetValue("new name", false)
	defObjects.SetValue([]int64{4, 5, 6}, false)
}

func TestModelFieldFieldsForceSetReadOnly(t *testing.T) {
	var m = &TestModelFields{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, false)
	)

	defID.SetValue(2, true)
	defName.SetValue("new name", true)
	defObjects.SetValue([]int64{4, 5, 6}, true)

	if m.ID != 2 {
		t.Errorf("expected %d, got %d", 2, m.ID)
	}

	if m.Name != "new name" {
		t.Errorf("expected %q, got %q", "new name", m.Name)
	}

	if m.Objects[0] != 4 {
		t.Errorf("expected %d, got %d", 4, m.Objects[0])
	}
}

func TestModelFormFields(t *testing.T) {
	var m = &TestModelFields{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, true)
	)

	var (
		formfieldID      = defID.FormField()
		formfieldName    = defName.FormField()
		formfieldObjects = defObjects.FormField()
	)

	if v, ok := formfieldID.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "ID" {
			t.Errorf("expected %q, got %q", "ID", v.Name())
		}

		if _, ok := v.Widget().(*widgets.NumberWidget[int]); !ok {
			t.Errorf("expected %t, got %t", true, ok)
		}
	}

	if v, ok := formfieldName.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "Name" {
			t.Errorf("expected %q, got %q", "Name", v.Name())
		}

		if _, ok := v.Widget().(*widgets.BaseWidget); !ok {
			t.Errorf("expected %t, got %t", true, ok)
		}
	}

	if v, ok := formfieldObjects.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "Objects" {
			t.Errorf("expected %q, got %q", "Objects", v.Name())
		}

		if _, ok := v.Widget().(*widgets.BaseWidget); !ok {
			t.Errorf("expected %t, got %t", true, ok)
		}
	}
}

func TestModelFormFieldsCustomType(t *testing.T) {
	var m = &TestModelFields{
		ID:      1,
		Name:    "name",
		Objects: []int64{1, 2, 3},
	}

	var (
		defID      = attrs.NewField(m, "ID", false, false, true)
		defName    = attrs.NewField(m, "Name", false, false, true)
		defObjects = attrs.NewField(m, "Objects", false, false, true)
	)

	goldcrest.Register(
		attrs.HookFormFieldForType,
		0, attrs.TypeGetter(func(f attrs.Field, t reflect.Type, v reflect.Value) (fields.Field, bool) {
			if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Int64 {
				var newF = fields.NewField(fields.S("text"))
				newF.FormWidget = &customTestWidget{widgets.NewBaseWidget(
					widgets.S("custom"),
					"", nil,
				)}
				return newF, true
			}
			return nil, false
		}),
	)

	var (
		formfieldID      = defID.FormField()
		formfieldName    = defName.FormField()
		formfieldObjects = defObjects.FormField()
	)

	if v, ok := formfieldID.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "ID" {
			t.Errorf("expected %q, got %q", "ID", v.Name())
		}

		if _, ok := v.Widget().(*widgets.NumberWidget[int]); !ok {
			t.Errorf("expected %t, got %t", true, ok)
		}
	}

	if v, ok := formfieldName.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "Name" {
			t.Errorf("expected %q, got %q", "Name", v.Name())
		}

		if _, ok := v.Widget().(*widgets.BaseWidget); !ok {
			t.Errorf("expected %t, got %t", true, ok)
		}
	}

	if v, ok := formfieldObjects.(*fields.BaseField); !ok {
		t.Errorf("expected %t, got %t", true, ok)
	} else {
		if v.Name() != "Objects" {
			t.Errorf("expected %q, got %q", "Objects", v.Name())
		}

		if _, ok := v.Widget().(*customTestWidget); !ok {
			t.Errorf("expected %t, got %t (%T)", true, ok, v.Widget())
		}
	}
}
