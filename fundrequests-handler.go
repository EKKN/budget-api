package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateFundRequestsRequest(reqBody *FundRequests) error {
	if reqBody.BudgetPostsID <= 0 {
		return fmt.Errorf("budget posts id must be filled")
	} else if reqBody.Date.IsZero() {
		return fmt.Errorf("date must be filled")
	} else if reqBody.Type == "" {
		return fmt.Errorf("type must be filled")
	} else if reqBody.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	} else if reqBody.Status == "" {
		return fmt.Errorf("status must be filled")
	}
	return nil
}

func (s *APIServer) checkFundRequestsForeignKey(primaryKey *PrimaryKeyID) (string, error) {

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: primaryKey.BudgetPostsID,
	}

	primaryKeyID, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKey)
	if err != nil {
		return "error DB", err
	}

	if primaryKeyID.BudgetPostsID == 0 {
		return "data budgetposts not found", nil
	}
	return "ok", nil
}

func (s *APIServer) GetAllFundRequests(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	fundRequest, err := s.Storage.FundRequestsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, fundRequest)

}

func (s *APIServer) GetFundRequestByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	fundRequest, err := s.Storage.FundRequestsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, fundRequest)

}

func (s *APIServer) CreateFundRequest(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &FundRequests{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateFundRequestsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: reqBody.BudgetPostsID,
	}

	message, err := s.checkFundRequestsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	fundRequest, err := s.Storage.FundRequestsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, fundRequest)

}

func (s *APIServer) UpdateFundRequest(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &FundRequests{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateFundRequestsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: reqBody.BudgetPostsID,
	}

	message, err := s.checkFundRequestsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedFundRequest, err := s.Storage.FundRequestsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if updatedFundRequest == nil {
		return respondWithError(requestLog, "data fund request not found", err)
	}
	return respondWithSuccess(requestLog, updatedFundRequest)

}

func (s *APIServer) DeleteFundRequest(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedFundRequest, err := s.Storage.FundRequestsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if deletedFundRequest == nil {
		return respondWithError(requestLog, "data fund request not found", err)
	}
	return respondWithSuccess(requestLog, deletedFundRequest)

}
