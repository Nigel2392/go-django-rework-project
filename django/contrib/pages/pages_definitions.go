package pages

import (
	"context"
	"net/http"

	"github.com/Nigel2392/django/contrib/admin"
	"github.com/Nigel2392/django/contrib/pages/models"
	"github.com/Nigel2392/django/core/contenttypes"
)

type PageDefinition struct {
	*contenttypes.ContentTypeDefinition
	AddPanels               func(r *http.Request, page Page) []admin.Panel
	EditPanels              func(r *http.Request, page Page) []admin.Panel
	GetForID                func(ctx context.Context, ref models.PageNode, id int64) (Page, error)
	OnReferenceUpdate       func(ctx context.Context, ref models.PageNode, id int64) error
	OnReferenceBeforeDelete func(ctx context.Context, ref models.PageNode, id int64) error
}

func (p *PageDefinition) Label() string {
	if p.GetLabel != nil {
		return p.GetLabel()
	}
	return ""
}

func (p *PageDefinition) Description() string {
	if p.GetDescription != nil {
		return p.GetDescription()
	}
	return ""
}

func (p *PageDefinition) ContentType() contenttypes.ContentType {
	if p.ContentTypeDefinition == nil {
		return nil
	}
	return p.ContentTypeDefinition.ContentType()
}

func (p *PageDefinition) AppLabel() string {
	return p.ContentType().AppLabel()
}

func (p *PageDefinition) Model() string {
	return p.ContentType().Model()
}
