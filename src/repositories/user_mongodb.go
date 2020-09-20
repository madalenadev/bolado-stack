package repositories

import (
	"bolado-stack/libs/database"
	"bolado-stack/src/domain"
	"context"
)

// IUserMongoRepository interface for user mongo repository
type IUserMongoRepository interface {
	ReadOne(ctx context.Context, ID string) (*domain.User, error)
}

type userMongoRepositoryImpl struct{}

// NewUserMongoDBRepository create new mongo repository for user
func NewUserMongoDBRepository(db database.IMongo) IUserMongoRepository {
	return userMongoRepositoryImpl{}
}

func (umr userMongoRepositoryImpl) ReadOne(ctx context.Context, ID string) (*domain.User, error) {
	return nil, nil
}
