package expense

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

// ListExpense return slice of expense info with given filter
// if filter object were empty, default filter will be used
func (s *service) ListExpense(ctx context.Context, userID string, filter domain.ExpenseFilter) ([]domain.ExpenseInfoDTO, error) {
	if filter.Size == 0 {
		filter.Size = 10
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.StartDate == "" {
		now := time.Now().UTC()
		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		filter.StartDate = firstDayOfMonth.Format(time.RFC3339Nano)
	}
	if filter.EndDate == "" {
		now := time.Now().UTC()
		nextMonth := now.Month() + 1
		nextYear := now.Year()
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}

		firstDayOfNextMonth := time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, time.UTC)
		lastDayOfMonth := firstDayOfNextMonth.Add(time.Second - 1)
		filter.EndDate = lastDayOfMonth.Format(time.RFC3339Nano)
	}

	expenses, err := s.repository.ListExpense(ctx, userID, filter)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("list latest expense: ", err)
		return nil, err
	}

	if expenses == nil {
		expenses = []domain.ExpenseInfoDTO{}
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

	if categories == nil {
		categories = []domain.ExpenseCategory{}
	}

	return categories, err
}
