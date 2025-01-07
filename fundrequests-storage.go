package main

import (
	"database/sql"
	"fmt"
)

type FundRequestsStorage interface {
	Create(*FundRequests) (*FundRequests, error)
	Delete(int64) (*FundRequests, error)
	Update(int64, *FundRequests) (*FundRequests, error)
	UpdateActive(int64, *FundRequests) (*FundRequests, error)
	GetDataId(int64) (*FundRequests, error)
	GetData() ([]*FundRequests, error)
	GetIdByName(string) (*FundRequests, error)
}

type FundRequestsStore struct {
	mysql *MysqlDB
}

func NewFundRequestsStorage(db *MysqlDB) *FundRequestsStore {
	return &FundRequestsStore{
		mysql: db,
	}
}
func (m *FundRequestsStore) GetIdByName(name string) (*FundRequests, error) {
	query := `SELECT id, budget_posts_id, date, type, amount, status, created_at, updated_at FROM fund_requests WHERE name = ?`
	row := m.mysql.db.QueryRow(query, name)

	fundRequest := &FundRequests{}
	err := row.Scan(&fundRequest.ID, &fundRequest.BudgetPostsID, &fundRequest.Date, &fundRequest.Type, &fundRequest.Amount, &fundRequest.Status, &fundRequest.CreatedAt, &fundRequest.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get fund request by id: %w", err)
	}
	return fundRequest, nil
}

func (m *FundRequestsStore) GetData() ([]*FundRequests, error) {
	query := `SELECT id, budget_posts_id, date, type, amount, status, created_at, updated_at FROM fund_requests`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund requests: %w", err)
	}
	defer rows.Close()

	var fundRequests []*FundRequests
	for rows.Next() {
		fundRequest := &FundRequests{}
		err := rows.Scan(
			&fundRequest.ID,
			&fundRequest.BudgetPostsID,
			&fundRequest.Date,
			&fundRequest.Type,
			&fundRequest.Amount,
			&fundRequest.Status,
			&fundRequest.CreatedAt,
			&fundRequest.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fund request: %w", err)
		}
		fundRequests = append(fundRequests, fundRequest)
	}
	return fundRequests, nil
}

func (m *FundRequestsStore) GetDataId(id int64) (*FundRequests, error) {
	query := `SELECT id, budget_posts_id, date, type, amount, status, created_at, updated_at FROM fund_requests WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	fundRequest := &FundRequests{}
	err := row.Scan(&fundRequest.ID, &fundRequest.BudgetPostsID, &fundRequest.Date, &fundRequest.Type, &fundRequest.Amount, &fundRequest.Status, &fundRequest.CreatedAt, &fundRequest.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fund request by id: %w", err)
	}
	return fundRequest, nil
}

func (m *FundRequestsStore) Create(fundRequest *FundRequests) (*FundRequests, error) {

	query := `INSERT INTO fund_requests (budget_posts_id, date, type, amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, fundRequest.BudgetPostsID, fundRequest.Date, fundRequest.Type, fundRequest.Amount, fundRequest.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to insert fund request: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newFundRequest, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new fund request: %w", err)
	}

	return newFundRequest, nil
}

func (m *FundRequestsStore) Delete(id int64) (*FundRequests, error) {
	deletedFundRequest, _ := m.GetDataId(id)

	query := `DELETE FROM fund_requests WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete fund request: %w", err)
	}

	return deletedFundRequest, nil
}

func (m *FundRequestsStore) Update(id int64, fundRequest *FundRequests) (*FundRequests, error) {
	query := `UPDATE fund_requests SET budget_posts_id = ?, date = ?, type = ?, amount = ?, status = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, fundRequest.BudgetPostsID, fundRequest.Date, fundRequest.Type, fundRequest.Amount, fundRequest.Status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update fund request: %w", err)
	}

	updatedFundRequest, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated fund request: %w", err)
	}

	return updatedFundRequest, nil
}

func (m *FundRequestsStore) UpdateActive(id int64, fundRequest *FundRequests) (*FundRequests, error) {
	query := `UPDATE fund_requests SET status = ?, updated_at = now() WHERE id = ?`

	_, err := m.mysql.db.Exec(query, fundRequest.Status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update fund request: %w", err)
	}

	updatedFundRequest, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated fund request: %w", err)
	}

	return updatedFundRequest, nil
}
