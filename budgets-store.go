package main

import (
	"database/sql"
	"fmt"
)

type BudgetsStorage interface {
	Create(*Budgets) (*Budgets, error)
	Delete(int64) (*Budgets, error)
	Update(int64, *Budgets) (*Budgets, error)
	GetDataId(int64) (*Budgets, error)
	GetData() ([]*Budgets, error)
	UpdateApproved(int64, *Budgets) (*Budgets, error)
	GetIdByName(string) (*Budgets, error)
}

type BudgetsStore struct {
	mysql *MysqlDB
}

func NewBudgetsStorage(db *MysqlDB) *BudgetsStore {
	return &BudgetsStore{
		mysql: db,
	}
}

func (m *BudgetsStore) GetIdByName(name string) (*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets WHERE name = ?`
	row := m.mysql.db.QueryRow(query, name)

	budgets := &Budgets{}
	err := row.Scan(&budgets.ID, &budgets.Name, &budgets.Description, &budgets.Periode, &budgets.IsApproved, &budgets.UnitsID, &budgets.CreatedAt, &budgets.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get budget by name: %w", err)
	}
	return budgets, nil
}

func (m *BudgetsStore) GetData() ([]*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}
	defer rows.Close()

	var budgetsList []*Budgets
	for rows.Next() {
		budgets := &Budgets{}
		err := rows.Scan(
			&budgets.ID,
			&budgets.Name,
			&budgets.Description,
			&budgets.Periode,
			&budgets.IsApproved,
			&budgets.UnitsID,
			&budgets.CreatedAt,
			&budgets.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgetsList = append(budgetsList, budgets)
	}
	return budgetsList, nil
}

func (m *BudgetsStore) GetDataId(id int64) (*Budgets, error) {
	query := `SELECT id, name, description, periode, is_approved, units_id, created_at, updated_at FROM budgets WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	budgets := &Budgets{}
	err := row.Scan(&budgets.ID, &budgets.Name, &budgets.Description, &budgets.Periode, &budgets.IsApproved, &budgets.UnitsID, &budgets.CreatedAt, &budgets.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget by id: %w", err)
	}
	return budgets, nil
}

func (m *BudgetsStore) Create(budgets *Budgets) (*Budgets, error) {
	query := `INSERT INTO budgets (name, description, periode, is_approved, units_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgets.Name, budgets.Description, budgets.Periode, budgets.IsApproved, budgets.UnitsID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgets, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget: %w", err)
	}

	return newBudgets, nil
}

func (m *BudgetsStore) Delete(id int64) (*Budgets, error) {
	deletedBudgets, _ := m.GetDataId(id)

	query := `DELETE FROM budgets WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget: %w", err)
	}

	return deletedBudgets, nil
}

func (m *BudgetsStore) Update(id int64, budgets *Budgets) (*Budgets, error) {
	query := `UPDATE budgets SET name = ?, description = ?, periode = ?, is_approved = ?, units_id = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgets.Name, budgets.Description, budgets.Periode, budgets.IsApproved, budgets.UnitsID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	updatedBudgets, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget: %w", err)
	}

	return updatedBudgets, nil
}

func (m *BudgetsStore) UpdateApproved(id int64, budgets *Budgets) (*Budgets, error) {

	query := `UPDATE budgets SET is_approved = ?, updated_at = now() WHERE id = ?`
	AppLog(query, budgets.IsApproved, id)
	_, err := m.mysql.db.Exec(query, budgets.IsApproved, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget post approval status: %w", err)
	}

	updatedBudgetPosts, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget post: %w", err)
	}

	return updatedBudgetPosts, nil
}
