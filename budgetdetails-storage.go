package main

import (
	"database/sql"
	"fmt"
)

type BudgetDetailsStorage interface {
	Create(*BudgetDetails) (*BudgetDetails, error)
	Delete(int64) (*BudgetDetails, error)
	Update(int64, *BudgetDetails) (*BudgetDetails, error)
	GetById(int64) (*BudgetDetails, error)
	GetAll() ([]*BudgetDetails, error)
}

type BudgetDetailsStore struct {
	db *sql.DB
}

func NewBudgetDetailsStorage(db *sql.DB) *BudgetDetailsStore {
	return &BudgetDetailsStore{
		db: db,
	}
}

func (s *BudgetDetailsStore) GetAll() ([]*BudgetDetails, error) {
	query := `SELECT id, budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at FROM budget_details`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget details: %w", err)
	}
	defer rows.Close()

	var budgetDetailsList []*BudgetDetails
	for rows.Next() {
		budgetDetail := &BudgetDetails{}
		err := rows.Scan(
			&budgetDetail.ID,
			&budgetDetail.BudgetsID,
			&budgetDetail.ActivitiesID,
			&budgetDetail.Description,
			&budgetDetail.Target,
			&budgetDetail.Quantity,
			&budgetDetail.UnitValue,
			&budgetDetail.Total,
			&budgetDetail.Terms,
			&budgetDetail.CreatedAt,
			&budgetDetail.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget detail: %w", err)
		}
		budgetDetailsList = append(budgetDetailsList, budgetDetail)
	}
	return budgetDetailsList, nil
}

func (s *BudgetDetailsStore) GetById(id int64) (*BudgetDetails, error) {
	query := `SELECT id, budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at FROM budget_details WHERE id = ?`
	row := s.db.QueryRow(query, id)

	budgetDetail := &BudgetDetails{}
	err := row.Scan(
		&budgetDetail.ID,
		&budgetDetail.BudgetsID,
		&budgetDetail.ActivitiesID,
		&budgetDetail.Description,
		&budgetDetail.Target,
		&budgetDetail.Quantity,
		&budgetDetail.UnitValue,
		&budgetDetail.Total,
		&budgetDetail.Terms,
		&budgetDetail.CreatedAt,
		&budgetDetail.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget detail by id: %w", err)
	}
	return budgetDetail, nil
}

func (s *BudgetDetailsStore) Create(budgetDetail *BudgetDetails) (*BudgetDetails, error) {
	query := `INSERT INTO budget_details (budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, now(), now())`
	result, err := s.db.Exec(query, budgetDetail.BudgetsID, budgetDetail.ActivitiesID, budgetDetail.Description, budgetDetail.Target, budgetDetail.Quantity, budgetDetail.UnitValue, budgetDetail.Total, budgetDetail.Terms)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget detail: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetDetailsStore) Delete(id int64) (*BudgetDetails, error) {
	budgetDetail, _ := s.GetById(id)
	if budgetDetail == nil {
		return nil, fmt.Errorf("budget detail not found")
	}
	query := `DELETE FROM budget_details WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget detail: %w", err)
	}
	return budgetDetail, nil
}

func (s *BudgetDetailsStore) Update(id int64, budgetDetail *BudgetDetails) (*BudgetDetails, error) {
	query := `UPDATE budget_details SET budgets_id = ?, activities_id = ?, description = ?, target = ?, quantity = ?, unit_value = ?, total = ?, terms = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, budgetDetail.BudgetsID, budgetDetail.ActivitiesID, budgetDetail.Description, budgetDetail.Target, budgetDetail.Quantity, budgetDetail.UnitValue, budgetDetail.Total, budgetDetail.Terms, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget detail: %w", err)
	}

	return s.GetById(id)
}
