package main

import (
	"database/sql"
	"fmt"
)

type BudgetsStorage interface {
	Create(*Budgets) (*Budgets, error)
	Delete(int64) (*Budgets, error)
	Update(int64, *Budgets) (*Budgets, error)
	GetById(int64) (*Budgets, error)
	GetAll() ([]*Budgets, error)
	UpdateApproved(int64, *Budgets) (*Budgets, error)
	GetByName(string) (*Budgets, error)
}

type BudgetsStore struct {
	db *sql.DB
}

func NewBudgetsStorage(db *sql.DB) *BudgetsStore {
	return &BudgetsStore{
		db: db,
	}
}

func (s *BudgetsStore) GetByName(name string) (*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets WHERE name = ?`
	row := s.db.QueryRow(query, name)

	budget := &Budgets{}
	err := row.Scan(&budget.ID, &budget.Name, &budget.Description, &budget.Periode, &budget.IsApproved, &budget.UnitsID, &budget.CreatedAt, &budget.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget by name: %w", err)
	}
	return budget, nil
}

func (s *BudgetsStore) GetAll() ([]*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}
	defer rows.Close()

	var budgetsList []*Budgets
	for rows.Next() {
		budget := &Budgets{}
		err := rows.Scan(
			&budget.ID,
			&budget.Name,
			&budget.Description,
			&budget.Periode,
			&budget.IsApproved,
			&budget.UnitsID,
			&budget.CreatedAt,
			&budget.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgetsList = append(budgetsList, budget)
	}
	return budgetsList, nil
}

func (s *BudgetsStore) GetById(id int64) (*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets WHERE id = ?`
	row := s.db.QueryRow(query, id)

	budget := &Budgets{}
	err := row.Scan(&budget.ID, &budget.Name, &budget.Description, &budget.Periode, &budget.IsApproved, &budget.UnitsID, &budget.CreatedAt, &budget.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget by id: %w", err)
	}
	return budget, nil
}

func (s *BudgetsStore) Create(budget *Budgets) (*Budgets, error) {
	query := `INSERT INTO budgets (name, description, periode, is_approved, units_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := s.db.Exec(query, budget.Name, budget.Description, budget.Periode, budget.IsApproved, budget.UnitsID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetsStore) Delete(id int64) (*Budgets, error) {
	deletedBudget, _ := s.GetById(id)

	query := `DELETE FROM budgets WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget: %w", err)
	}

	return deletedBudget, nil
}

func (s *BudgetsStore) Update(id int64, budget *Budgets) (*Budgets, error) {
	query := `UPDATE budgets SET name = ?, description = ?, periode = ?, is_approved = ?, units_id = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, budget.Name, budget.Description, budget.Periode, budget.IsApproved, budget.UnitsID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return s.GetById(id)
}

func (s *BudgetsStore) UpdateApproved(id int64, budget *Budgets) (*Budgets, error) {
	query := `UPDATE budgets SET is_approved = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, budget.IsApproved, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget approval status: %w", err)
	}

	return s.GetById(id)
}
