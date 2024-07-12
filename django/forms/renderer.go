package forms

import (
	"embed"
	"io/fs"

	"github.com/Nigel2392/django/core/assert"
	"github.com/Nigel2392/django/core/filesystem"
	"github.com/Nigel2392/django/core/filesystem/tpl"
)

//go:embed assets/**
var formTemplates embed.FS

func init() {
	var templates, err = fs.Sub(formTemplates, "assets/templates")
	assert.True(err == nil, "failed to get form templates")

	tpl.Add(tpl.Config{
		AppName: "forms",
		FS:      templates,
		Bases:   []string{},
		Matches: filesystem.MatchAnd(
			filesystem.MatchPrefix("forms/widgets/"),
			filesystem.MatchOr(
				filesystem.MatchExt(".html"),
			),
		),
	})
}
