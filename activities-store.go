package main

import (
	"database/sql"
	"fmt"
)

type ActivitiesStorage interface {
	Create(*Activities) (*Activities, error)
	Delete(int64) (*Activities, error)
	Update(int64, *Activities) (*Activities, error)
	UpdateActive(int64, *Activities) (*Activities, error)
	GetDataId(int64) (*Activities, error)
	GetData() ([]*Activities, error)
	GetIdByName(string) (*Activities, error)
}

type ActivitiesStore struct {
	mysql *MysqlDB
}

func NewActivitiesStorage(db *MysqlDB) *ActivitiesStore {
	return &ActivitiesStore{
		mysql: db,
	}
}

func (m *ActivitiesStore) GetIdByName(name string) (*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities WHERE name = ?`
	row := m.mysql.db.QueryRow(query, name)

	activity := &Activities{}
	err := row.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.IsActive, &activity.CreatedAt, &activity.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get activity by id: %w", err)
	}
	return activity, nil
}

func (m *ActivitiesStore) GetData() ([]*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer rows.Close()

	var activities []*Activities
	for rows.Next() {
		activity := &Activities{}
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Description,
			&activity.IsActive,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (m *ActivitiesStore) GetDataId(id int64) (*Activities, error) {
	query := `SELECT id, name, description, is_active, created_at, updated_at FROM activities WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	activity := &Activities{}
	err := row.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.IsActive, &activity.CreatedAt, &activity.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get activity by id: %w", err)
	}
	return activity, nil
}

func (m *ActivitiesStore) Create(activity *Activities) (*Activities, error) {

	query := `INSERT INTO activities (name, description, is_active, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, activity.Name, activity.Description, activity.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to insert activity: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newActivity, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new activity: %w", err)
	}

	return newActivity, nil
}

func (m *ActivitiesStore) Delete(id int64) (*Activities, error) {
	deletedActivity, _ := m.GetDataId(id)

	query := `DELETE FROM activities WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete activity: %w", err)
	}

	return deletedActivity, nil
}

func (m *ActivitiesStore) Update(id int64, activity *Activities) (*Activities, error) {
	query := `UPDATE activities SET name = ?, description = ?, is_active = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, activity.Name, activity.Description, activity.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	updatedActivity, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated activity: %w", err)
	}

	return updatedActivity, nil
}

func (m *ActivitiesStore) UpdateActive(id int64, activity *Activities) (*Activities, error) {
	query := `UPDATE activities SET  is_active = ?, updated_at = now() WHERE id = ?`

	_, err := m.mysql.db.Exec(query, activity.IsActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	updatedActivity, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated activity: %w", err)
	}

	return updatedActivity, nil
}
