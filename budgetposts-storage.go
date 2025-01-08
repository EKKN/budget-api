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
	GetById(int64) (*BudgetPosts, error)
	GetAll() ([]*BudgetPosts, error)
	GetByName(string) (*BudgetPosts, error)
}

type BudgetPostsStore struct {
	db *sql.DB
}

func NewBudgetPostsStorage(db *sql.DB) *BudgetPostsStore {
	return &BudgetPostsStore{db: db}
}

func (s *BudgetPostsStore) GetByName(name string) (*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts WHERE name = ?`
	row := s.db.QueryRow(query, name)

	budgetPost := &BudgetPosts{}
	err := row.Scan(&budgetPost.ID, &budgetPost.Name, &budgetPost.Description, &budgetPost.IsActive, &budgetPost.CreatedAt, &budgetPost.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget post by name: %w", err)
	}
	return budgetPost, nil
}

func (s *BudgetPostsStore) GetAll() ([]*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget posts: %w", err)
	}
	defer rows.Close()

	var budgetPostsList []*BudgetPosts
	for rows.Next() {
		budgetPost := &BudgetPosts{}
		err := rows.Scan(&budgetPost.ID, &budgetPost.Name, &budgetPost.Description, &budgetPost.IsActive, &budgetPost.CreatedAt, &budgetPost.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget post: %w", err)
		}
		budgetPostsList = append(budgetPostsList, budgetPost)
	}
	return budgetPostsList, nil
}

func (s *BudgetPostsStore) GetById(id int64) (*BudgetPosts, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM budget_posts WHERE id = ?`
	row := s.db.QueryRow(query, id)

	budgetPost := &BudgetPosts{}
	err := row.Scan(&budgetPost.ID, &budgetPost.Name, &budgetPost.Description, &budgetPost.IsActive, &budgetPost.CreatedAt, &budgetPost.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget post by id: %w", err)
	}
	return budgetPost, nil
}

func (s *BudgetPostsStore) Create(budgetPost *BudgetPosts) (*BudgetPosts, error) {
	existingBudgetPost, err := s.GetByName(budgetPost.Name)
	if err != nil {
		return nil, fmt.Errorf("error checking name: %w", err)
	}
	if existingBudgetPost != nil {
		return existingBudgetPost, fmt.Errorf("name already in use")
	}

	query := `INSERT INTO budget_posts (name, description, is_active, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := s.db.Exec(query, budgetPost.Name, budgetPost.Description, budgetPost.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget post: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetPostsStore) Delete(id int64) (*BudgetPosts, error) {
	budgetPost, _ := s.GetById(id)
	if budgetPost == nil {
		return nil, fmt.Errorf("budget post not found")
	}

	query := `DELETE FROM budget_posts WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget post: %w", err)
	}

	return budgetPost, nil
}

func (s *BudgetPostsStore) Update(id int64, budgetPost *BudgetPosts) (*BudgetPosts, error) {
	existingBudgetPost, err := s.GetByName(budgetPost.Name)
	if err != nil {
		return nil, fmt.Errorf("error checking name: %w", err)
	}
	if existingBudgetPost != nil && existingBudgetPost.ID != id {
		return existingBudgetPost, fmt.Errorf("name already in use")
	}

	query := `UPDATE budget_posts SET name = ?, description = ?, is_active = ?, updated_at = now() WHERE id = ?`
	_, err = s.db.Exec(query, budgetPost.Name, budgetPost.Description, budgetPost.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget post: %w", err)
	}

	return s.GetById(id)
}

func (s *BudgetPostsStore) UpdateActive(id int64, budgetPost *BudgetPosts) (*BudgetPosts, error) {
	query := `UPDATE budget_posts SET is_active = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, budgetPost.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget post: %w", err)
	}

	return s.GetById(id)
}
