// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package list

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/Nigel2392/django/core/attrs"
	"strings"
)

type ColumnGroup[T attrs.Definer] struct {
	Definitons attrs.Definitions
	Columns    []ListColumn[T]
	Instance   T
}

func (c *ColumnGroup[T]) AddColumn(column ListColumn[T]) {
	c.Columns = append(c.Columns, column)
}

func (c *ColumnGroup[T]) Render() string {
	var component = c.Component()
	var b strings.Builder
	var ctx = context.Background()
	component.Render(ctx, &b)
	return b.String()
}

func (c *ColumnGroup[T]) Component() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<tr class=\"column-group\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, column := range c.Columns {
			templ_7745c5c3_Err = column.Component(c.Definitons, c.Instance).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
