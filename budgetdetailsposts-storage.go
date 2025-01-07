package main

import (
	"database/sql"
	"fmt"
)

type BudgetDetailsPostsStorage interface {
	Create(*BudgetDetailsPosts) (*BudgetDetailsPosts, error)
	Delete(int64) (*BudgetDetailsPosts, error)
	Update(int64, *BudgetDetailsPosts) (*BudgetDetailsPosts, error)
	GetDataByID(int64) (*BudgetDetailsPosts, error)
	GetData() ([]*BudgetDetailsPosts, error)
}

type BudgetDetailsPostsStore struct {
	mysql *MysqlDB
}

func NewBudgetDetailsPostsStorage(db *MysqlDB) *BudgetDetailsPostsStore {
	return &BudgetDetailsPostsStore{
		mysql: db,
	}
}

func (m *BudgetDetailsPostsStore) GetData() ([]*BudgetDetailsPosts, error) {
	query := `SELECT id, budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at FROM budget_details_posts`
	rows, err := m.mysql.db.Query(query)
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

func (m *BudgetDetailsPostsStore) GetDataByID(id int64) (*BudgetDetailsPosts, error) {
	query := `SELECT id, budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at FROM budget_details_posts WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

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

func (m *BudgetDetailsPostsStore) Create(budgetDetailsPost *BudgetDetailsPosts) (*BudgetDetailsPosts, error) {
	query := `INSERT INTO budget_details_posts (budget_details_id, budget_posts_id, planned_amount, approved_amount, usage_amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgetDetailsPost.BudgetDetailsID, budgetDetailsPost.BudgetPostsID, budgetDetailsPost.PlannedAmount, budgetDetailsPost.ApprovedAmount, budgetDetailsPost.UsageAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget details post: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgetDetailsPost, err := m.GetDataByID(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget details post: %w", err)
	}

	return newBudgetDetailsPost, nil
}

func (m *BudgetDetailsPostsStore) Delete(id int64) (*BudgetDetailsPosts, error) {
	deletedBudgetDetailsPost, _ := m.GetDataByID(id)

	query := `DELETE FROM budget_details_posts WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget details post: %w", err)
	}

	return deletedBudgetDetailsPost, nil
}

func (m *BudgetDetailsPostsStore) Update(id int64, budgetDetailsPost *BudgetDetailsPosts) (*BudgetDetailsPosts, error) {
	query := `UPDATE budget_details_posts SET budget_details_id = ?, budget_posts_id = ?, planned_amount = ?, approved_amount = ?, usage_amount = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetDetailsPost.BudgetDetailsID, budgetDetailsPost.BudgetPostsID, budgetDetailsPost.PlannedAmount, budgetDetailsPost.ApprovedAmount, budgetDetailsPost.UsageAmount, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget details post: %w", err)
	}

	updatedBudgetDetailsPost, err := m.GetDataByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget details post: %w", err)
	}

	return updatedBudgetDetailsPost, nil
}
