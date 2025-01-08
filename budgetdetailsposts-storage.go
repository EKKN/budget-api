package main

import (
	"database/sql"
	"fmt"
)

type BudgetDetailsPostsStorage interface {
	Create(*BudgetDetailsPosts) (*BudgetDetailsPosts, error)
	Delete(int64) (*BudgetDetailsPosts, error)
	Update(int64, *BudgetDetailsPosts) (*BudgetDetailsPosts, error)
	GetById(int64) (*BudgetDetailsPosts, error)
	GetAll() ([]*BudgetDetailsPosts, error)
}

type BudgetDetailsPostsStore struct {
	db *sql.DB
}

func NewBudgetDetailsPostsStorage(db *sql.DB) *BudgetDetailsPostsStore {
	return &BudgetDetailsPostsStore{
		db: db,
	}
}

func (s *BudgetDetailsPostsStore) GetAll() ([]*BudgetDetailsPosts, error) {
	query := `SELECT id, budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at FROM budget_details_posts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget details posts: %w", err)
	}
	defer rows.Close()

	var budgetDetailsPostsList []*BudgetDetailsPosts
	for rows.Next() {
		budgetDetailsPost := &BudgetDetailsPosts{}
		err := rows.Scan(
			&budgetDetailsPost.ID,
			&budgetDetailsPost.BudgetDetailsID,
			&budgetDetailsPost.BudgetPostsID,
			&budgetDetailsPost.PlannedAmount,
			&budgetDetailsPost.ApprovedAmount,
			&budgetDetailsPost.UsageAmount,
			&budgetDetailsPost.CreatedAt,
			&budgetDetailsPost.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget details post: %w", err)
		}
		budgetDetailsPostsList = append(budgetDetailsPostsList, budgetDetailsPost)
	}
	return budgetDetailsPostsList, nil
}

func (s *BudgetDetailsPostsStore) GetById(id int64) (*BudgetDetailsPosts, error) {
	query := `SELECT id, budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at FROM budget_details_posts WHERE id = ?`
	row := s.db.QueryRow(query, id)

	budgetDetailsPost := &BudgetDetailsPosts{}
	err := row.Scan(
		&budgetDetailsPost.ID,
		&budgetDetailsPost.BudgetDetailsID,
		&budgetDetailsPost.BudgetPostsID,
		&budgetDetailsPost.PlannedAmount,
		&budgetDetailsPost.ApprovedAmount,
		&budgetDetailsPost.UsageAmount,
		&budgetDetailsPost.CreatedAt,
		&budgetDetailsPost.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget details post by id: %w", err)
	}
	return budgetDetailsPost, nil
}

func (s *BudgetDetailsPostsStore) Create(post *BudgetDetailsPosts) (*BudgetDetailsPosts, error) {
	query := `INSERT INTO budget_details_posts (budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := s.db.Exec(query, post.BudgetDetailsID, post.BudgetPostsID, post.PlannedAmount, post.ApprovedAmount, post.UsageAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget details post: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetDetailsPostsStore) Delete(id int64) (*BudgetDetailsPosts, error) {
	post, _ := s.GetById(id)
	if post == nil {
		return nil, fmt.Errorf("budget details post not found")
	}
	query := `DELETE FROM budget_details_posts WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget details post: %w", err)
	}
	return post, nil
}

func (s *BudgetDetailsPostsStore) Update(id int64, post *BudgetDetailsPosts) (*BudgetDetailsPosts, error) {
	query := `UPDATE budget_details_posts SET budget_details_id = ?, budget_posts_id = ?, planned_amount = ?, approved_amount = ?, usage_amount = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, post.BudgetDetailsID, post.BudgetPostsID, post.PlannedAmount, post.ApprovedAmount, post.UsageAmount, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget details post: %w", err)
	}
	return s.GetById(id)
}
