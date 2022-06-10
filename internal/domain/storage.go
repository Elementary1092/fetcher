package domain

import "context"

type Storage interface {
    SaveAll(ctx context.Context, posts []*Post) error
}
