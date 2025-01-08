package main

import (
	"database/sql"
	"fmt"
	"time"
)

type ActivitiesStorage interface {
	Create(*Activities) (*Activities, error)
	Delete(int64) (*Activities, error)
	Update(int64, *Activities) (*Activities, error)
	UpdateActive(int64, *Activities) (*Activities, error)
	GetById(int64) (*Activities, error)
	GetAll() ([]*Activities, error)
	GetByName(string) (*Activities, error)
}

type ActivitiesStore struct {
	db *sql.DB
}

func NewActivitiesStorage(db *sql.DB) *ActivitiesStore {
	return &ActivitiesStore{
		db: db,
	}
}

func scanActivity(row *sql.Row) (*Activities, error) {
	activity := &Activities{}
	err := row.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.IsActive, &activity.CreatedAt, &activity.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan activity: %w", err)
	}
	return activity, nil
}

func scanActivities(rows *sql.Rows) ([]*Activities, error) {
	var activities []*Activities
	for rows.Next() {
		activity := &Activities{}
		err := rows.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.IsActive, &activity.CreatedAt, &activity.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (s *ActivitiesStore) GetByName(name string) (*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities WHERE name = ?`
	row := s.db.QueryRow(query, name)
	return scanActivity(row)
}

func (s *ActivitiesStore) GetAll() ([]*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer rows.Close()
	return scanActivities(rows)
}

func (s *ActivitiesStore) GetById(id int64) (*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities WHERE id = ?`
	row := s.db.QueryRow(query, id)
	return scanActivity(row)
}

func (s *ActivitiesStore) Create(activity *Activities) (*Activities, error) {
	query := `INSERT INTO activities (name, description, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, activity.Name, activity.Description, activity.IsActive, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert activity: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return s.GetById(lastInsertID)
}

func (s *ActivitiesStore) Delete(id int64) (*Activities, error) {
	activity, _ := s.GetById(id)
	if activity == nil {
		return nil, fmt.Errorf("activity not found")
	}
	query := `DELETE FROM activities WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete activity: %w", err)
	}
	return activity, nil
}

func (s *ActivitiesStore) Update(id int64, activity *Activities) (*Activities, error) {
	query := `UPDATE activities SET name = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, activity.Name, activity.Description, activity.IsActive, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}
	return s.GetById(id)
}

func (s *ActivitiesStore) UpdateActive(id int64, activity *Activities) (*Activities, error) {
	query := `UPDATE activities SET is_active = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, activity.IsActive, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}
	return s.GetById(id)
}
