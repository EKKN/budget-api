package main

import (
	"database/sql"
	"fmt"
)

type FundRequestDetailsStorage interface {
	Create(*FundRequestDetails) (*FundRequestDetails, error)
	Delete(int64) (*FundRequestDetails, error)
	Update(int64, *FundRequestDetails) (*FundRequestDetails, error)
	GetDataId(int64) (*FundRequestDetails, error)
	GetData() ([]*FundRequestDetails, error)
}

type FundRequestDetailsStore struct {
	mysql *MysqlDB
}

func NewFundRequestDetailsStorage(db *MysqlDB) *FundRequestDetailsStore {
	return &FundRequestDetailsStore{
		mysql: db,
	}
}
func (m *FundRequestDetailsStore) GetData() ([]*FundRequestDetails, error) {
	query := `SELECT id, fund_requests_id, activities_id, budget_details_id, amount, recommendation, created_at, updated_at FROM fund_request_details`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund request details: %w", err)
	}
	defer rows.Close()

	var fundRequestDetails []*FundRequestDetails
	for rows.Next() {
		fundRequestDetail := &FundRequestDetails{}
		err := rows.Scan(
			&fundRequestDetail.ID,
			&fundRequestDetail.FundRequestsID,
			&fundRequestDetail.ActivitiesID,
			&fundRequestDetail.BudgetDetailsID,
			&fundRequestDetail.Amount,
			&fundRequestDetail.Recommendation,
			&fundRequestDetail.CreatedAt,
			&fundRequestDetail.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fund request detail: %w", err)
		}
		fundRequestDetails = append(fundRequestDetails, fundRequestDetail)
	}
	return fundRequestDetails, nil
}

func (m *FundRequestDetailsStore) GetDataId(id int64) (*FundRequestDetails, error) {
	query := `SELECT id, fund_requests_id, activities_id, budget_details_id, amount, recommendation, created_at, updated_at FROM fund_request_details WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	fundRequestDetail := &FundRequestDetails{}
	err := row.Scan(&fundRequestDetail.ID, &fundRequestDetail.FundRequestsID, &fundRequestDetail.ActivitiesID, &fundRequestDetail.BudgetDetailsID, &fundRequestDetail.Amount, &fundRequestDetail.Recommendation, &fundRequestDetail.CreatedAt, &fundRequestDetail.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fund request detail by id: %w", err)
	}
	return fundRequestDetail, nil
}

func (m *FundRequestDetailsStore) Create(fundRequestDetail *FundRequestDetails) (*FundRequestDetails, error) {

	query := `INSERT INTO fund_request_details (fund_requests_id, activities_id, budget_details_id, amount, recommendation, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, fundRequestDetail.FundRequestsID, fundRequestDetail.ActivitiesID, fundRequestDetail.BudgetDetailsID, fundRequestDetail.Amount, fundRequestDetail.Recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to insert fund request detail: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newFundRequestDetail, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new fund request detail: %w", err)
	}

	return newFundRequestDetail, nil
}

func (m *FundRequestDetailsStore) Delete(id int64) (*FundRequestDetails, error) {
	deletedFundRequestDetail, _ := m.GetDataId(id)

	query := `DELETE FROM fund_request_details WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete fund request detail: %w", err)
	}

	return deletedFundRequestDetail, nil
}

func (m *FundRequestDetailsStore) Update(id int64, fundRequestDetail *FundRequestDetails) (*FundRequestDetails, error) {
	query := `UPDATE fund_request_details SET fund_requests_id = ?, activities_id = ?, budget_details_id = ?, amount = ?, recommendation = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, fundRequestDetail.FundRequestsID, fundRequestDetail.ActivitiesID, fundRequestDetail.BudgetDetailsID, fundRequestDetail.Amount, fundRequestDetail.Recommendation, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update fund request detail: %w", err)
	}

	updatedFundRequestDetail, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated fund request detail: %w", err)
	}

	return updatedFundRequestDetail, nil
}
