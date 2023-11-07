package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
	"github.com/ryanadiputraa/spendr-backend/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	log        logger.Logger
	repository domain.UserRepository
}

func NewService(log logger.Logger, validator validator.Validator, repository domain.UserRepository) domain.AuthService {
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

func (s *service) Signin(ctx context.Context, email, password string) (*domain.User, error) {
	logPrefix := "signin user: "
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewError(domain.BadRequest, "no user registered with given email")
		}
		s.log.Error(logPrefix, err)
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.NewError(domain.Unauthorized, "password didn't match")
	}

	return user, nil
}
