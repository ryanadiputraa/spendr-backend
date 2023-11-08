package user

import (
	"context"

	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

type service struct {
	log        logger.Logger
	repository domain.UserRepository
}

func NewService(log logger.Logger, repository domain.UserRepository) domain.UserService {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) GetUserData(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repository.FindUserByID(ctx, userID)
	if err != nil {
		s.log.Error("get user data: ", err)
		return nil, err
	}
	return user, nil
}

func (s *service) ListSupportedCurrency(ctx context.Context) []string {
	currencies := make([]string, 0)
	for c := range domain.Currencies {
		currencies = append(currencies, c)
	}
	return currencies
}
