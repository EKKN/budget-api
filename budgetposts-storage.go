package main

import (
	"database/sql"
	"fmt"
)

type BudgetPostsStorage interface {
	Create(*BudgetPosts) (*BudgetPosts, error)
	Delete(int64) (*BudgetPosts, error)
	Update(int64, *BudgetPosts) (*BudgetPosts, error)
	UpdateActive(int64, *BudgetPosts) (*BudgetPosts, error)
	GetDataId(int64) (*BudgetPosts, error)
	GetData() ([]*BudgetPosts, error)
	GetIdByName(string) (*BudgetPosts, error)
}

type BudgetPostsStore struct {
	mysql *MysqlDB
}

func NewBudgetPostsStorage(db *MysqlDB) *BudgetPostsStore {
	return &BudgetPostsStore{
		mysql: db,
	}
}

func (m *BudgetPostsStore) GetIdByName(name string) (*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts WHERE name = ?`
	row := m.mysql.db.QueryRow(query, name)

	budgetPosts := &BudgetPosts{}
	err := row.Scan(&budgetPosts.ID, &budgetPosts.Name, &budgetPosts.Description, &budgetPosts.IsActive, &budgetPosts.CreatedAt, &budgetPosts.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get budget post by name: %w", err)
	}
	return budgetPosts, nil
}

func (m *BudgetPostsStore) GetData() ([]*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget posts: %w", err)
	}
	defer rows.Close()

	var budgetPostsList []*BudgetPosts
	for rows.Next() {
		budgetPosts := &BudgetPosts{}
		err := rows.Scan(
			&budgetPosts.ID,
			&budgetPosts.Name,
			&budgetPosts.Description,
			&budgetPosts.IsActive,
			&budgetPosts.CreatedAt,
			&budgetPosts.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget post: %w", err)
		}
		budgetPostsList = append(budgetPostsList, budgetPosts)
	}
	return budgetPostsList, nil
}

func (m *BudgetPostsStore) GetDataId(id int64) (*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	budgetPosts := &BudgetPosts{}
	err := row.Scan(&budgetPosts.ID, &budgetPosts.Name, &budgetPosts.Description, &budgetPosts.IsActive, &budgetPosts.CreatedAt, &budgetPosts.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget post by id: %w", err)
	}
	return budgetPosts, nil
}

func (m *BudgetPostsStore) Create(budgetPosts *BudgetPosts) (*BudgetPosts, error) {
	budgetPostsByName, err := m.GetIdByName(budgetPosts.Name)
	if err != nil {
		return nil, fmt.Errorf("error checking name: %w", err)
	}
	if budgetPostsByName != nil {
		return budgetPostsByName, fmt.Errorf("name already in use")
	}

	query := `INSERT INTO budget_posts (name, description, is_active, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgetPosts.Name, budgetPosts.Description, budgetPosts.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget post: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgetPosts, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget post: %w", err)
	}

	return newBudgetPosts, nil
}

func (m *BudgetPostsStore) Delete(id int64) (*BudgetPosts, error) {
	deletedBudgetPosts, _ := m.GetDataId(id)

	query := `DELETE FROM budget_posts WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget post: %w", err)
	}

	return deletedBudgetPosts, nil
}

func (m *BudgetPostsStore) Update(id int64, budgetPosts *BudgetPosts) (*BudgetPosts, error) {
	budgetPostsByName, err := m.GetIdByName(budgetPosts.Name)
	if err != nil {
		return nil, fmt.Errorf("error checking name: %w", err)
	}
	if budgetPostsByName != nil && budgetPostsByName.ID != id {
		return budgetPostsByName, fmt.Errorf("name already in use")
	}

	query := `UPDATE budget_posts SET name = ?, description = ?, is_active = ?, updated_at = now() WHERE id = ?`
	_, err = m.mysql.db.Exec(query, budgetPosts.Name, budgetPosts.Description, budgetPosts.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget post: %w", err)
	}

	updatedBudgetPosts, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget post: %w", err)
	}

	return updatedBudgetPosts, nil
}

func (m *BudgetPostsStore) UpdateActive(id int64, budgetPosts *BudgetPosts) (*BudgetPosts, error) {
	query := `UPDATE budget_posts SET is_active = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetPosts.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget post: %w", err)
	}

	updatedBudgetPosts, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget post: %w", err)
	}

	return updatedBudgetPosts, nil
}
