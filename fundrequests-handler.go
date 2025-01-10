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

// newPrimaryKey := &PrimaryKeyID{
// 	FundRequestDetailsID: id,
// 	ActivitiesID:         reqBody.ActivitiesID,
// 	FundRequestsID:       reqBody.FundRequestsID,
// 	BudgetDetailsID:      reqBody.BudgetDetailsID,
// }

// message, err := s.validateFundRequestDetailsForeignKey(newPrimaryKey, true)
//
//	if err != nil {
//		return respondWithError(requestLog, message, err)
//	}

func (s *APIServer) validateFundRequestsForeignKey(primaryKey *PrimaryKeyID, validateSelfID bool, checkSelfOnly bool) (string, error) {

	newKey := &PrimaryKeyID{
		FundRequestsID: primaryKey.FundRequestsID,
		BudgetPostsID:  primaryKey.BudgetPostsID,
	}

	storedKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newKey)
	if err != nil {
		return "database error", err
	}

	if validateSelfID && storedKey.FundRequestsID == 0 {
		return "fund requests not found", fmt.Errorf("fund requests not found")
	}
	if !checkSelfOnly {

		if storedKey.BudgetPostsID == 0 {
			return "data budget post not found", fmt.Errorf("data budget post not found")
		}
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

	message, err := s.validateFundRequestsForeignKey(newPrimaryKey, false, false)
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
		FundRequestsID: id,
		BudgetPostsID:  reqBody.BudgetPostsID,
	}

	message, err := s.validateFundRequestsForeignKey(newPrimaryKey, true, false)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedFundRequest, err := s.Storage.FundRequestsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	// if updatedFundRequest == nil {
	// 	return respondWithError(requestLog, "data fund request not found", err)
	// }
	return respondWithSuccess(requestLog, updatedFundRequest)

}

func (s *APIServer) DeleteFundRequest(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	newPrimaryKey := &PrimaryKeyID{
		FundRequestsID: id,
	}

	message, err := s.validateFundRequestsForeignKey(newPrimaryKey, true, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	deletedFundRequest, err := s.Storage.FundRequestsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	// if deletedFundRequest == nil {
	// 	return respondWithError(requestLog, "data fund request not found", err)
	// }
	return respondWithSuccess(requestLog, deletedFundRequest)

}
