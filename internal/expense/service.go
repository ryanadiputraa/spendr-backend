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

func (s *service) AddExpense(ctx context.Context, userID string, dto domain.ExpenseDTO) (*domain.Expense, error) {
	expense, err := domain.NewExpense(
		uuid.NewString(),
		userID,
		dto,
	)
	if err != nil {
		return nil, domain.NewError(domain.BadRequest, err.Error())
	}

	if err := s.repository.AddExpense(ctx, *expense); err != nil {
		s.log.Error("add expense: ", err)
		return nil, err
	}

	return expense, nil
}

func (s *service) ListLatestExpense(ctx context.Context, userID string, limit int) ([]domain.ExpenseInfoDTO, error) {
	if limit == 0 {
		limit = 10
	}

	expenses, err := s.repository.ListLatestExpense(ctx, userID, limit)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("list latest expense: ", err)
		return nil, err
	}

	return expenses, nil
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
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("list expense category: ", err)
		return nil, err
	}

	return categories, err
}
