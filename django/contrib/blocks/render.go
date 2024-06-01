package blocks

import (
	"context"
	"io"

	"github.com/Nigel2392/django/core/ctx"
)

func RenderBlockForm(w io.Writer, widget *BlockWidget, ctx *BlockContext, errors []error) error {
	return RenderBlockWidget(w, widget, ctx, errors).Render(context.Background(), w)
}

func RenderBlock(w io.Writer, block Block, value any, context ctx.Context) error {
	return block.Render(w, value, context)
}
