package main

import (
	"database/sql"
	"fmt"
)

// BudgetDetailsStorage interface
type BudgetDetailsStorage interface {
	Create(*BudgetDetails) (*BudgetDetails, error)
	Delete(int64) (*BudgetDetails, error)
	Update(int64, *BudgetDetails) (*BudgetDetails, error)
	GetDataByID(int64) (*BudgetDetails, error)
	GetData() ([]*BudgetDetails, error)
}

// BudgetDetailsStore struct
type BudgetDetailsStore struct {
	mysql *MysqlDB
}

// NewBudgetDetailsStorage initializes BudgetDetailsStore
func NewBudgetDetailsStorage(db *MysqlDB) *BudgetDetailsStore {
	return &BudgetDetailsStore{
		mysql: db,
	}
}

// GetData retrieves all budget details
func (m *BudgetDetailsStore) GetData() ([]*BudgetDetails, error) {
	query := `SELECT id, budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at FROM budget_details`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget details: %w", err)
	}
	defer rows.Close()

	var budgetDetails []*BudgetDetails
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
		budgetDetails = append(budgetDetails, budgetDetail)
	}
	return budgetDetails, nil
}

// GetDataByID retrieves budget detail by ID
func (m *BudgetDetailsStore) GetDataByID(id int64) (*BudgetDetails, error) {
	query := `SELECT id, budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at FROM budget_details WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

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

// Create inserts a new budget detail
func (m *BudgetDetailsStore) Create(budgetDetail *BudgetDetails) (*BudgetDetails, error) {
	query := `INSERT INTO budget_details (budgets_id, activities_id, description, target, quantity, unit_value, total, terms, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgetDetail.BudgetsID, budgetDetail.ActivitiesID, budgetDetail.Description, budgetDetail.Target, budgetDetail.Quantity, budgetDetail.UnitValue, budgetDetail.Total, budgetDetail.Terms)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget detail: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgetDetail, err := m.GetDataByID(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget detail: %w", err)
	}

	return newBudgetDetail, nil
}

// Delete removes a budget detail by ID
func (m *BudgetDetailsStore) Delete(id int64) (*BudgetDetails, error) {
	deletedBudgetDetail, _ := m.GetDataByID(id)

	query := `DELETE FROM budget_details WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget detail: %w", err)
	}

	return deletedBudgetDetail, nil
}

// Update modifies a budget detail
func (m *BudgetDetailsStore) Update(id int64, budgetDetail *BudgetDetails) (*BudgetDetails, error) {
	query := `UPDATE budget_details SET budgets_id = ?, activities_id = ?, description = ?, target = ?, quantity = ?, unit_value = ?, total = ?, terms = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetDetail.BudgetsID, budgetDetail.ActivitiesID, budgetDetail.Description, budgetDetail.Target, budgetDetail.Quantity, budgetDetail.UnitValue, budgetDetail.Total, budgetDetail.Terms, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget detail: %w", err)
	}
	AppLog(1)

	updatedBudgetDetail, err := m.GetDataByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget detail: %w", err)
	}

	return updatedBudgetDetail, nil
}
