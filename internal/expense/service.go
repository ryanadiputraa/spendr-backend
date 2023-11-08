package expense

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ryanadiputraa/spendr-backend/internal/domain"
	"github.com/ryanadiputraa/spendr-backend/pkg/logger"
)

type service struct {
	log        logger.Logger
	repository domain.ExpenseRepository
}

func NewService(log logger.Logger, repository domain.ExpenseRepository) domain.ExpenseService {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) AddExpenseCategory(ctx context.Context, userID string, dto domain.ExpenseCategoryDTO) (*domain.ExpenseCategory, error) {
	category := domain.NewExpenseCategory(
		uuid.NewString(),
		dto.Category,
		dto.Ico,
		userID,
	)

	if err := s.repository.AddExpenseCategory(ctx, *category); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == domain.PQErrDuplicate {
			return nil, domain.NewError(domain.BadRequest, "category already exists")
		}
		s.log.Error("add expense category: ", err)
		return nil, err
	}

	return category, nil
}

func (s *service) ListExpenseCategory(ctx context.Context, userID string) ([]domain.ExpenseCategory, error) {
	categories, err := s.repository.ListExpenseCategory(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.Error("list expense category: ", err)
			return nil, err
		}
	}

	return categories, err
}
