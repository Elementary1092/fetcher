package api

import (
	"context"
	"encoding/json"
	"github.com/elem1092/fetcher/internal/domain"
	"github.com/elem1092/fetcher/pkg/logging"
	"io"
	"net/http"
	"strconv"
	"time"
)

type service struct {
	logger  *logging.Logger
	storage domain.Storage
	baseUrl string
}

func NewService(logger *logging.Logger, storage domain.Storage, url string) domain.Service {
	return &service{
		logger:  logger,
		storage: storage,
		baseUrl: url,
	}
}

func (s *service) FetchPage(ctx context.Context, page int32) error {
	url := s.baseUrl + "?page=" + strconv.Itoa(int(page))
	s.logger.Infof("Fetching %s", url)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		s.logger.Warnf("Error while parsing %s: %v", url, err)
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		s.logger.Warnf("Error while parsing body of %s: %v", url, err)
		return err
	}

	s.logger.Info("Unmarshalling json")
	posts := make([]*domain.PostDTO, 0)
	err = json.Unmarshal(body, posts)
	if err != nil {
		s.logger.Errorf("Got error while unmrshalling: %v", err)
		return err
	}

	s.logger.Infof("Converting into sotrable posts")
	storablePosts := make([]*domain.Post, len(posts))
	for i, postDTO := range posts {
		storablePosts[i] = postDTO.FormPost()
	}

	s.logger.Info("Trying to save posts into the database")
	return s.storage.SaveAll(ctx, storablePosts)
}
