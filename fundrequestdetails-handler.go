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

// newPrimaryKey := &PrimaryKeyID{
// 	BudgetsID: id,
// }

// message, err := s.validateActivitiesForeignKey(newPrimaryKey, true)
// if err != nil {
// 	return respondWithError(requestLog, message, err)
// }

func (s *APIServer) validateFundRequestDetailsForeignKey(primaryKey *PrimaryKeyID, validateSelfID bool, checkSelfOnly bool) (string, error) {

	newKey := &PrimaryKeyID{
		FundRequestDetailsID: primaryKey.FundRequestDetailsID,
		ActivitiesID:         primaryKey.ActivitiesID,
		FundRequestsID:       primaryKey.FundRequestsID,
		BudgetDetailsID:      primaryKey.BudgetDetailsID,
	}

	storedKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newKey)
	if err != nil {
		return "database error", err
	}

	if validateSelfID && storedKey.FundRequestDetailsID == 0 {
		return "fund reqeust details not found", fmt.Errorf("fund reqeust details not found")
	}

	if !checkSelfOnly {
		if storedKey.ActivitiesID == 0 {
			return "activities not found", fmt.Errorf("activities not found")
		}

		if storedKey.FundRequestsID == 0 {
			return "fund request not found", fmt.Errorf("fund request not found")
		}

		if storedKey.BudgetDetailsID == 0 {
			return "budget details not found", fmt.Errorf("budget details not found")
		}

	}

	return "ok", nil
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

	newPrimaryKey := &PrimaryKeyID{
		ActivitiesID:    reqBody.ActivitiesID,
		FundRequestsID:  reqBody.FundRequestsID,
		BudgetDetailsID: reqBody.BudgetDetailsID,
	}

	message, err := s.validateFundRequestDetailsForeignKey(newPrimaryKey, false, false)
	if err != nil {
		return respondWithError(requestLog, message, err)
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

	newPrimaryKey := &PrimaryKeyID{
		FundRequestDetailsID: id,
		ActivitiesID:         reqBody.ActivitiesID,
		FundRequestsID:       reqBody.FundRequestsID,
		BudgetDetailsID:      reqBody.BudgetDetailsID,
	}

	message, err := s.validateFundRequestDetailsForeignKey(newPrimaryKey, true, false)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedFundRequestDetail, err := s.Storage.FundRequestDetailsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	// if updatedFundRequestDetail == nil {
	// 	return respondWithError(requestLog, "data fund request not found", err)
	// }
	return respondWithSuccess(requestLog, updatedFundRequestDetail)

}

func (s *APIServer) DeleteFundRequestDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	newPrimaryKey := &PrimaryKeyID{
		FundRequestDetailsID: id,
	}

	message, err := s.validateFundRequestDetailsForeignKey(newPrimaryKey, true, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	deletedFundRequestDetail, err := s.Storage.FundRequestDetailsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	// if deletedFundRequestDetail == nil {
	// 	return respondWithError(requestLog, "data fund request not found", err)
	// }
	return respondWithSuccess(requestLog, deletedFundRequestDetail)

}
