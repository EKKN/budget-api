package main

import (
	"database/sql"
	"fmt"
)

type BudgetDetailsPostsRecommendationsStorage interface {
	Create(*BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error)
	Delete(int64) (*BudgetDetailsPostsRecommendations, error)
	Update(int64, *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error)
	GetDataId(int64) (*BudgetDetailsPostsRecommendations, error)
	GetData() ([]*BudgetDetailsPostsRecommendations, error)
}

type BudgetDetailsPostsRecommendationsStore struct {
	mysql *MysqlDB
}

func NewBudgetDetailsPostsRecommendationsStorage(db *MysqlDB) *BudgetDetailsPostsRecommendationsStore {
	return &BudgetDetailsPostsRecommendationsStore{
		mysql: db,
	}
}
func (m *BudgetDetailsPostsRecommendationsStore) GetData() ([]*BudgetDetailsPostsRecommendations, error) {
	query := `SELECT id, budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at FROM budget_details_posts_recommendations`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget details posts recommendations: %w", err)
	}
	defer rows.Close()

	var budgetDetailsPostsRecommendations []*BudgetDetailsPostsRecommendations
	for rows.Next() {
		budgetDetailsPostsRecommendation := &BudgetDetailsPostsRecommendations{}
		err := rows.Scan(
			&budgetDetailsPostsRecommendation.ID,
			&budgetDetailsPostsRecommendation.BudgetDetailsPostsID,
			&budgetDetailsPostsRecommendation.UserGroupsID,
			&budgetDetailsPostsRecommendation.Recommendation,
			&budgetDetailsPostsRecommendation.CreatedAt,
			&budgetDetailsPostsRecommendation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget details posts recommendation: %w", err)
		}
		budgetDetailsPostsRecommendations = append(budgetDetailsPostsRecommendations, budgetDetailsPostsRecommendation)
	}
	return budgetDetailsPostsRecommendations, nil
}

func (m *BudgetDetailsPostsRecommendationsStore) GetDataId(id int64) (*BudgetDetailsPostsRecommendations, error) {
	query := `SELECT id, budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at FROM budget_details_posts_recommendations WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	budgetDetailsPostsRecommendation := &BudgetDetailsPostsRecommendations{}
	err := row.Scan(&budgetDetailsPostsRecommendation.ID, &budgetDetailsPostsRecommendation.BudgetDetailsPostsID, &budgetDetailsPostsRecommendation.UserGroupsID, &budgetDetailsPostsRecommendation.Recommendation, &budgetDetailsPostsRecommendation.CreatedAt, &budgetDetailsPostsRecommendation.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget details posts recommendation by id: %w", err)
	}
	return budgetDetailsPostsRecommendation, nil
}

func (m *BudgetDetailsPostsRecommendationsStore) Create(budgetDetailsPostsRecommendation *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error) {

	query := `INSERT INTO budget_details_posts_recommendations (budget_details_posts_id, user_groups_id, recommendation, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgetDetailsPostsRecommendation.BudgetDetailsPostsID, budgetDetailsPostsRecommendation.UserGroupsID, budgetDetailsPostsRecommendation.Recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget details posts recommendation: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgetDetailsPostsRecommendation, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget details posts recommendation: %w", err)
	}

	return newBudgetDetailsPostsRecommendation, nil
}

func (m *BudgetDetailsPostsRecommendationsStore) Delete(id int64) (*BudgetDetailsPostsRecommendations, error) {
	deletedBudgetDetailsPostsRecommendation, _ := m.GetDataId(id)

	query := `DELETE FROM budget_details_posts_recommendations WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget details posts recommendation: %w", err)
	}

	return deletedBudgetDetailsPostsRecommendation, nil
}

func (m *BudgetDetailsPostsRecommendationsStore) Update(id int64, budgetDetailsPostsRecommendation *BudgetDetailsPostsRecommendations) (*BudgetDetailsPostsRecommendations, error) {
	query := `UPDATE budget_details_posts_recommendations SET budget_details_posts_id = ?, user_groups_id = ?, recommendation = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetDetailsPostsRecommendation.BudgetDetailsPostsID, budgetDetailsPostsRecommendation.UserGroupsID, budgetDetailsPostsRecommendation.Recommendation, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget details posts recommendation: %w", err)
	}

	updatedBudgetDetailsPostsRecommendation, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget details posts recommendation: %w", err)
	}

	return updatedBudgetDetailsPostsRecommendation, nil
}
