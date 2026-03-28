package auth

import (
	"context"

	"portlyn/internal/domain"
)

type contextKey string

const userContextKey contextKey = "user"
const groupIDsContextKey contextKey = "group_ids"

func ContextWithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func ContextWithGroupIDs(ctx context.Context, groupIDs []uint) context.Context {
	return context.WithValue(ctx, groupIDsContextKey, append([]uint(nil), groupIDs...))
}

func UserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}

func GroupIDsFromContext(ctx context.Context) ([]uint, bool) {
	groupIDs, ok := ctx.Value(groupIDsContextKey).([]uint)
	if !ok {
		return nil, false
	}
	return append([]uint(nil), groupIDs...), true
}
