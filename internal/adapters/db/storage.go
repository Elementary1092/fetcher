package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/elem1092/fetcher/internal/domain"
	"github.com/elem1092/fetcher/pkg/client/postgre"
	"github.com/elem1092/fetcher/pkg/logging"
	"github.com/jackc/pgconn"
	"strings"
)

type postgreSQLStorage struct {
	db     postgre.Client
	logger *logging.Logger
}

func NewPostgreSQLStorage(db postgre.Client, logger *logging.Logger) domain.Storage {
	return &postgreSQLStorage{
		db:     db,
		logger: logger,
	}
}

func (p *postgreSQLStorage) SaveAll(ctx context.Context, posts []*domain.Post) error {
	queryTail, err := p.formFieldsToBeInserted(posts)
	if err != nil {
		return err
	}
	p.logger.Info("Creating and executing SQL query")

	query := `INSERT INTO posts (user_id, title, body) VALUES` + queryTail

	_, err = p.db.Exec(ctx, query)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			p.logger.Errorf(
				"Got SQL error: %s; %s; where: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where,
			)

			return pgErr
		}

		p.logger.Errorf("Caught: %v", err)
		return err
	}

	p.logger.Info("Successfully saved record into the database")
	return nil
}

func (p *postgreSQLStorage) formFieldsToBeInserted(posts []*domain.Post) (string, error) {
	var str strings.Builder
	for _, post := range posts {
		if str.Len() != 0 {
			if _, err := str.WriteString(", "); err != nil {
				p.logger.Warnf("Got error while inserting field: %v", err)
				return "", err
			}
		}

		if _, err := str.WriteString(
			fmt.Sprintf("(%d, %s, %s", post.UserId, post.Title, post.Title)); err != nil {
			p.logger.Warnf("Got error while inserting field: %v", err)
			return "", err
		}
	}

	return str.String(), nil
}
