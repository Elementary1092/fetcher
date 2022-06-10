package api

import (
	"context"
	"github.com/elem1092/fetcher/internal/domain"
	fetch "github.com/elem1092/fetcher/pkg/client/grpc"
	"github.com/elem1092/fetcher/pkg/logging"
	"sync"
)

type server struct {
	service   domain.Service
	logger    *logging.Logger
	mutex     *sync.Mutex
	pagesLeft int32
	err       error
	fetch.UnimplementedFetchServiceServer
}

func NewServer(service domain.Service, logger *logging.Logger, mutex sync.Mutex, pages int32) fetch.FetchServiceServer {
	return &server{
		service:   service,
		logger:    logger,
		mutex:     &mutex,
		pagesLeft: pages,
	}
}

func (s *server) StartFetching(ctx context.Context, request *fetch.FetchRequest) (*fetch.FetchStatus, error) {
	s.logger.Info("starting to fetch")
	if request.GetPages() != 0 {
		s.pagesLeft = request.GetPages()
	}

	pages := s.pagesLeft
	for i := int32(1); i <= pages; i++ {
		go func(page int32) {
			err := s.service.FetchPage(ctx, page)
			if err != nil {
				s.mutex.Lock()
				s.err = err
				s.mutex.Unlock()

				return
			}

			s.mutex.Lock()
			s.pagesLeft--
			s.mutex.Unlock()
		}(i)
	}

	return nil, nil
}

func (s *server) GetStatus(ctx context.Context, r *fetch.EmptyMessage) (*fetch.FetchStatus, error) {
	s.logger.Info("handling GetStatus request")

	s.mutex.Lock()
	err := s.err
	s.mutex.Unlock()
	if err != nil {
		return &fetch.FetchStatus{StatusCode: 2}, err
	}

	s.mutex.Lock()
	pages := s.pagesLeft
	s.mutex.Unlock()
	if pages == 0 {
		return &fetch.FetchStatus{StatusCode: 1}, nil
	}

	return &fetch.FetchStatus{StatusCode: 0}, nil
}

func (s *server) GetError(ctx context.Context, r *fetch.EmptyMessage) (*fetch.EmptyMessage, error) {
	s.logger.Info("handling GetError request")

	s.mutex.Lock()
	err := s.err
	s.mutex.Unlock()

	return nil, err
}
