package restdb

import "context"

type ctxKeyPrefix struct{}

func WithPrefix(ctx context.Context, prefix string) context.Context {
	return context.WithValue(ctx, ctxKeyPrefix{}, prefix)
}

func FromPrefix(ctx context.Context) string {
	if s, ok := ctx.Value(ctxKeyPrefix{}).(string); ok {
		return s
	}
	return ""
}
