package main

import (
	"database/sql"
	"fmt"
)

type PrimaryKeyIDStorage interface {
	GetPrimaryKey(PrimaryKeyID *PrimaryKeyID) (*PrimaryKeyID, error)
}

type PrimaryKeyIDStore struct {
	db *sql.DB
}

func NewPrimaryKeyIDStorage(db *sql.DB) *PrimaryKeyIDStore {
	return &PrimaryKeyIDStore{
		db: db,
	}
}
func (s *PrimaryKeyIDStore) GetPrimaryKey(primaryKey *PrimaryKeyID) (*PrimaryKeyID, error) {
	queries := map[string]*int64{
		"SELECT id FROM budgets WHERE id = ?":              &primaryKey.BudgetsID,
		"SELECT id FROM budget_posts WHERE id = ?":         &primaryKey.BudgetPostsID,
		"SELECT id FROM activities WHERE id = ?":           &primaryKey.ActivitiesID,
		"SELECT id FROM budget_details WHERE id = ?":       &primaryKey.BudgetDetailsID,
		"SELECT id FROM budget_details_posts WHERE id = ?": &primaryKey.BudgetDetailsPostsID,
		"SELECT id FROM fund_requests WHERE id = ?":        &primaryKey.FundRequestsID,
	}

	for query, idPtr := range queries {
		if *idPtr > 0 {
			row := s.db.QueryRow(query, *idPtr)
			var id int64
			if err := row.Scan(&id); err != nil {
				if err == sql.ErrNoRows {
					*idPtr = 0
				} else {
					return nil, fmt.Errorf("failed to get ID: %w", err)
				}
			}
			*idPtr = id
		}
	}

	return primaryKey, nil
}
