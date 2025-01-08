package main

import (
	"database/sql"
	"fmt"
)

type BudgetDetailPostRecStorage interface {
	Create(*BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error)
	Delete(int64) (*BudgetDetailsPostsRecommendations, error)
	Update(int64, *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error)
	GetById(int64) (*BudgetDetailsPostsRecommendations, error)
	GetAll() ([]*BudgetDetailsPostsRecommendations, error)
}

type BudgetDetailPostRecStore struct {
	db *sql.DB
}

func NewBudgetDetailPostRecStorage(db *sql.DB) *BudgetDetailPostRecStore {
	return &BudgetDetailPostRecStore{
		db: db,
	}
}

func (s *BudgetDetailPostRecStore) GetAll() ([]*BudgetDetailsPostsRecommendations, error) {
	query := `SELECT id, budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at FROM budget_details_posts_recommendations`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget detail post recommendations: %w", err)
	}
	defer rows.Close()

	var recs []*BudgetDetailsPostsRecommendations
	for rows.Next() {
		rec := &BudgetDetailsPostsRecommendations{}
		err := rows.Scan(
			&rec.ID,
			&rec.BudgetDetailsPostsID,
			&rec.UserGroupsID,
			&rec.Recommendation,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget detail post recommendation: %w", err)
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

func (s *BudgetDetailPostRecStore) GetById(id int64) (*BudgetDetailsPostsRecommendations, error) {
	query := `SELECT id, budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at FROM budget_details_posts_recommendations WHERE id = ?`
	row := s.db.QueryRow(query, id)

	rec := &BudgetDetailsPostsRecommendations{}
	err := row.Scan(
		&rec.ID,
		&rec.BudgetDetailsPostsID,
		&rec.UserGroupsID,
		&rec.Recommendation,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget detail post recommendation by id: %w", err)
	}
	return rec, nil
}

func (s *BudgetDetailPostRecStore) Create(rec *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error) {
	query := `INSERT INTO budget_details_posts_recommendations (budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := s.db.Exec(query, rec.BudgetDetailsPostsID, rec.UserGroupsID, rec.Recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget detail post recommendation: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *BudgetDetailPostRecStore) Delete(id int64) (*BudgetDetailsPostsRecommendations, error) {
	rec, _ := s.GetById(id)
	if rec == nil {
		return nil, fmt.Errorf("budget detail post recommendation not found")
	}
	query := `DELETE FROM budget_details_posts_recommendations WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget detail post recommendation: %w", err)
	}
	return rec, nil
}

func (s *BudgetDetailPostRecStore) Update(id int64, rec *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error) {
	query := `UPDATE budget_details_posts_recommendations SET budget_details_posts_id = ?, user_groups_id = ?, recommendation = ?, updated_at = now() WHERE id = ?`
	_, err := s.db.Exec(query, rec.BudgetDetailsPostsID, rec.UserGroupsID, rec.Recommendation, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget detail post recommendation: %w", err)
	}
	return s.GetById(id)
}
