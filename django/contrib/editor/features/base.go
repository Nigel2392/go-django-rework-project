package features

import (
	"context"
	"fmt"
	"io"
	"maps"

	"github.com/Nigel2392/django/contrib/editor"
	"github.com/Nigel2392/django/core/ctx"
	"github.com/Nigel2392/django/forms/media"
)

var _ editor.BaseFeature = (*BaseFeature)(nil)
var _ editor.FeatureBlockRenderer = (*BaseFeature)(nil)

type BaseFeature struct {
	Type          string
	Extra         map[string]interface{}
	JSConstructor string
	JSFiles       []string
	CSSFles       []string
	RenderBlock   func(*FeatureBlock) *FeatureBlock
}

// Name returns the name of the feature.
func (b *BaseFeature) Name() string {
	return b.Type
}

// Config returns the configuration of the feature.
func (b *BaseFeature) Config(widgetContext ctx.Context) map[string]interface{} {
	var config = make(map[string]interface{})
	maps.Copy(config, b.Extra)
	return config
}

// Constructor returns the JS class name of the feature.
func (b *BaseFeature) Constructor() string {
	return b.JSConstructor
}

// Media return's the feature's static / media files.
func (b *BaseFeature) Media() media.Media {
	var m = media.NewMedia()
	for _, js := range b.JSFiles {
		m.AddJS(&media.JSAsset{URL: js})
	}
	for _, css := range b.CSSFles {
		m.AddCSS(&media.CSSAsset{URL: css})
	}
	return m
}

// Render returns the rendered HTML of the feature.
func (b *BaseFeature) Render(d editor.BlockData) editor.FeatureBlock {
	var block = &FeatureBlock{
		FeatureObject: b,
		FeatureData:   d,
		FeatureName:   d.Type,
		Identifier:    d.ID,
	}
	if b.RenderBlock != nil {
		block = b.RenderBlock(block)
	}
	return block
}

type Block struct {
	BaseFeature
	RenderFunc func(b editor.FeatureBlock, c context.Context, w io.Writer) error
}

func (b *Block) RenderBlock(fb editor.FeatureBlock, c context.Context, w io.Writer) error {
	if b.RenderFunc != nil {
		return b.RenderFunc(fb, c, w)
	}
	panic(
		fmt.Sprintf("feature '%s' does not implement RenderBlock method", b.Type),
	)
}
