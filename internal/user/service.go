package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
)

type service struct {
	log        logger.Logger
	repository domain.UserRepository
}

func NewService(log logger.Logger, validator validator.Validator, repository domain.UserRepository) domain.UserService {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) Signup(ctx context.Context, dto domain.UserDTO) (*domain.User, error) {
	logPrefix := "signup user: "
	user, err := domain.NewUser(
		uuid.NewString(),
		dto.Email,
		dto.Password,
		dto.FirstName,
		dto.LastName,
		dto.Picture,
		dto.Currency,
	)
	if err != nil {
		s.log.Warn(logPrefix, err)
		return nil, domain.NewError(domain.BadRequest, err.Error())
	}

	if err := user.HashPassword(); err != nil {
		s.log.Warn(logPrefix, err)
		return nil, domain.NewError(domain.BadRequest, err.Error())
	}

	if err = s.repository.AddUser(ctx, user); err != nil {
		s.log.Warn(logPrefix, err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == domain.PQErrDuplicate {
				return nil, domain.NewError(domain.BadRequest, "email already registered")
			}
		}
		return nil, domain.NewError(domain.BadRequest, err.Error())
	}

	return user, nil
}
