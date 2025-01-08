package main

import (
	"database/sql"
	"fmt"
	"time"
)

type BudgetCapsStorage interface {
	Create(*BudgetCaps) (*BudgetCaps, error)
	Delete(int64) (*BudgetCaps, error)
	Update(int64, *BudgetCaps) (*BudgetCaps, error)
	UpdateAmount(int64, *BudgetCaps) (*BudgetCaps, error)
	GetById(int64) (*BudgetCaps, error)
	GetAll() ([]*BudgetCaps, error)
}

type BudgetCapsStore struct {
	db *sql.DB
}

func NewBudgetCapsStorage(db *sql.DB) *BudgetCapsStore {
	return &BudgetCapsStore{
		db: db,
	}
}

func scanBudgetCap(row *sql.Row) (*BudgetCaps, error) {
	budgetCap := &BudgetCaps{}
	err := row.Scan(&budgetCap.ID, &budgetCap.BudgetsID, &budgetCap.BudgetPostsID, &budgetCap.Amount, &budgetCap.CreatedAt, &budgetCap.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan budget cap: %w", err)
	}
	return budgetCap, nil
}

func scanBudgetCaps(rows *sql.Rows) ([]*BudgetCaps, error) {
	var budgetCapsList []*BudgetCaps
	for rows.Next() {
		budgetCap := &BudgetCaps{}
		err := rows.Scan(&budgetCap.ID, &budgetCap.BudgetsID, &budgetCap.BudgetPostsID, &budgetCap.Amount, &budgetCap.CreatedAt, &budgetCap.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget cap: %w", err)
		}
		budgetCapsList = append(budgetCapsList, budgetCap)
	}
	return budgetCapsList, nil
}

func (s *BudgetCapsStore) GetAll() ([]*BudgetCaps, error) {
	query := `SELECT id, budgets_id, budget_posts_id, amount, created_at, updated_at FROM budget_caps`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget caps: %w", err)
	}
	defer rows.Close()
	return scanBudgetCaps(rows)
}

func (s *BudgetCapsStore) GetById(id int64) (*BudgetCaps, error) {
	query := `SELECT id, budgets_id, budget_posts_id, amount, created_at, updated_at FROM budget_caps WHERE id = ?`
	row := s.db.QueryRow(query, id)
	return scanBudgetCap(row)
}

func (s *BudgetCapsStore) Create(budgetCap *BudgetCaps) (*BudgetCaps, error) {
	query := `INSERT INTO budget_caps (budgets_id, budget_posts_id, amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, budgetCap.BudgetsID, budgetCap.BudgetPostsID, budgetCap.Amount, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget cap: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetCapsStore) Delete(id int64) (*BudgetCaps, error) {
	budgetCap, _ := s.GetById(id)
	if budgetCap == nil {
		return nil, fmt.Errorf("budget cap not found")
	}
	query := `DELETE FROM budget_caps WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget cap: %w", err)
	}
	return budgetCap, nil
}

func (s *BudgetCapsStore) Update(id int64, budgetCap *BudgetCaps) (*BudgetCaps, error) {
	query := `UPDATE budget_caps SET budgets_id = ?, budget_posts_id = ?, amount = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, budgetCap.BudgetsID, budgetCap.BudgetPostsID, budgetCap.Amount, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget cap: %w", err)
	}
	return s.GetById(id)
}

func (s *BudgetCapsStore) UpdateAmount(id int64, budgetCap *BudgetCaps) (*BudgetCaps, error) {
	query := `UPDATE budget_caps SET amount = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, budgetCap.Amount, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget cap amount: %w", err)
	}
	return s.GetById(id)
}
