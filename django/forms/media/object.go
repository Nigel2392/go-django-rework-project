package media

import (
	"html/template"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

type MediaObject struct {
	Css *orderedmap.OrderedMap[string, Asset]
	Js  *orderedmap.OrderedMap[string, Asset]
}

func NewMedia() *MediaObject {
	return &MediaObject{
		Css: orderedmap.NewOrderedMap[string, Asset](),
		Js:  orderedmap.NewOrderedMap[string, Asset](),
	}
}

func (m *MediaObject) String() string {
	var ret strings.Builder
	ret.WriteString("CSS:")
	if m.Css != nil && m.Css.Len() > 0 {
		ret.WriteString("\n")
	}
	for head := m.Css.Front(); head != nil; head = head.Next() {
		ret.WriteString("  ")
		ret.WriteString(head.Value.String())
		ret.WriteString("\n")
	}

	ret.WriteString("JS:")
	if m.Js != nil && m.Js.Len() > 0 {
		ret.WriteString("\n")
	}
	for head := m.Js.Front(); head != nil; head = head.Next() {
		ret.WriteString("  ")
		ret.WriteString(head.Value.String())
		ret.WriteString("\n")
	}
	return ret.String()
}

func (m *MediaObject) Merge(other Media) Media {
	if other == nil {
		return m
	}

	var (
		otherCss = other.CSSList()
		otherJs  = other.JSList()
	)

	for _, v := range otherCss {
		var strV = v.String()
		if _, ok := m.Css.Get(strV); !ok {
			m.Css.Set(strV, v)
		}
	}

	for _, v := range otherJs {
		var strV = v.String()
		if _, ok := m.Js.Get(strV); !ok {
			m.Js.Set(strV, v)
		}
	}

	return m
}

func (m *MediaObject) AddJS(list ...Asset) {
	for _, js := range list {
		m.Js.Set(js.String(), js)
	}
}

func (m *MediaObject) AddCSS(list ...Asset) {
	for _, css := range list {
		m.Css.Set(css.String(), css)
	}
}

func (m *MediaObject) JS() []template.HTML {
	var ret = make([]template.HTML, 0, m.Js.Len())
	for head := m.Js.Front(); head != nil; head = head.Next() {
		ret = append(ret, head.Value.Render())
	}
	return ret
}

func (m *MediaObject) CSS() []template.HTML {
	var ret = make([]template.HTML, 0, m.Css.Len())
	for head := m.Css.Front(); head != nil; head = head.Next() {
		ret = append(ret, head.Value.Render())
	}
	return ret
}

func (m *MediaObject) JSList() []Asset {
	var ret = make([]Asset, 0, m.Js.Len())
	for head := m.Js.Front(); head != nil; head = head.Next() {
		ret = append(ret, head.Value)
	}
	return ret
}

func (m *MediaObject) CSSList() []Asset {
	var ret = make([]Asset, 0, m.Css.Len())
	for head := m.Css.Front(); head != nil; head = head.Next() {
		ret = append(ret, head.Value)
	}
	return ret
}
