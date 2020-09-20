package services

import (
	"bolado-stack/src/domain"
	"context"

	"github.com/estrategiahq/questions-madalena/repositories"
)

// IUserService is the interface for user mongo repository
type IUserService interface {
	ReadOne(ctx context.Context, ID string) (*domain.User, error)
}

type userServiceImpl struct {
	repositories repositories.Container
}

// NewUserService create a new service for user
func NewUserService(repositories repositories.Container) IUserService {
	return userServiceImpl{repositories}
}

func (us userServiceImpl) ReadOne(ctx context.Context, ID string) (*domain.User, error) {
	return nil, nil
}
