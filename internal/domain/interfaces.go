package domain

import "context"

type Service interface {
	FetchPage(ctx context.Context, page int32) error
}
