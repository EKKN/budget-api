package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateFundRequestDetailsRequest(reqBody *FundRequestDetails) error {
	if reqBody.FundRequestsID <= 0 {
		return fmt.Errorf("fund requests id must be filled")
	} else if reqBody.ActivitiesID <= 0 {
		return fmt.Errorf("activities id must be filled")
	} else if reqBody.BudgetDetailsID <= 0 {
		return fmt.Errorf("budget details id must be filled")
	} else if reqBody.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	} else if reqBody.Recommendation == "" {
		return fmt.Errorf("recommendation must be filled")
	}
	return nil
}

func (s *APIServer) GetAllFundRequestDetails(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	fundRequestDetails, err := s.Storage.FundRequestDetailsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, fundRequestDetails)

}

func (s *APIServer) GetFundRequestDetailByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	fundRequestDetail, err := s.Storage.FundRequestDetailsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, fundRequestDetail)

}

func (s *APIServer) CreateFundRequestDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &FundRequestDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateFundRequestDetailsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKeyID := &PrimaryKeyID{
		ActivitiesID:    reqBody.ActivitiesID,
		FundRequestsID:  reqBody.FundRequestsID,
		BudgetDetailsID: reqBody.BudgetDetailsID,
	}

	primaryKeyID, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKeyID)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if primaryKeyID.ActivitiesID == 0 {
		return respondWithError(requestLog, "data activities not found", fmt.Errorf("data activities not found"))
	}

	if primaryKeyID.FundRequestsID == 0 {
		return respondWithError(requestLog, "data fund request not found", fmt.Errorf("data fund request not found"))
	}

	if primaryKeyID.BudgetDetailsID == 0 {
		return respondWithError(requestLog, "data budget details not found", fmt.Errorf("data budget details not found"))
	}

	fundRequestDetail, err := s.Storage.FundRequestDetailsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, fundRequestDetail)

}

func (s *APIServer) UpdateFundRequestDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &FundRequestDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateFundRequestDetailsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKeyID := &PrimaryKeyID{
		ActivitiesID:    reqBody.ActivitiesID,
		FundRequestsID:  reqBody.FundRequestsID,
		BudgetDetailsID: reqBody.BudgetDetailsID,
	}

	primaryKeyID, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKeyID)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if primaryKeyID.ActivitiesID == 0 {
		return respondWithError(requestLog, "data activities not found", fmt.Errorf("data activities not found"))
	}

	if primaryKeyID.FundRequestsID == 0 {
		return respondWithError(requestLog, "data fund request not found", fmt.Errorf("data fund request not found"))
	}

	if primaryKeyID.BudgetDetailsID == 0 {
		return respondWithError(requestLog, "data budget details not found", fmt.Errorf("data budget details not found"))
	}

	updatedFundRequestDetail, err := s.Storage.FundRequestDetailsStorage.Update(id, reqBody)
	if err != nil || updatedFundRequestDetail == nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, updatedFundRequestDetail)

}

func (s *APIServer) DeleteFundRequestDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedFundRequestDetail, err := s.Storage.FundRequestDetailsStorage.Delete(id)
	if err != nil || deletedFundRequestDetail == nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, deletedFundRequestDetail)

}
