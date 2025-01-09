package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetDetailsRequest(reqBody *BudgetDetails) error {
	if reqBody.BudgetsID <= 0 {
		return fmt.Errorf("budgets id must be filled")
	} else if reqBody.ActivitiesID <= 0 {
		return fmt.Errorf("activities id must be filled")
	} else if reqBody.Description == "" {
		return fmt.Errorf("description must be filled")
	} else if reqBody.Target.IsZero() {
		return fmt.Errorf("target must be filled")
	} else if reqBody.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	} else if reqBody.UnitValue <= 0 {
		return fmt.Errorf("unit value must be greater than 0")
	} else if reqBody.Total <= 0 {
		return fmt.Errorf("total must be greater than 0")
	} else if reqBody.Terms <= 0 {
		return fmt.Errorf("terms must be greater than 0")
	}
	return nil
}

func (s *APIServer) checkBudgetsDetailsForeignKey(primaryKey *PrimaryKeyID) (string, error) {

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:    primaryKey.BudgetsID,
		ActivitiesID: primaryKey.ActivitiesID,
	}

	primaryKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKey)
	if err != nil {
		return "error DB", err
	}

	if primaryKey.ActivitiesID == 0 {
		return "data activities not found", fmt.Errorf("data activities not found")
	}

	if primaryKey.BudgetsID == 0 {
		return "data budgets not found", fmt.Errorf("data budgets not found")
	}
	return "ok", nil
}

func (s *APIServer) GetAllBudgetDetails(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgetDetails, err := s.Storage.BudgetDetailsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetails)

}

func (s *APIServer) GetBudgetDetailByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budgetDetail, err := s.Storage.BudgetDetailsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetail)

}

func (s *APIServer) CreateBudgetDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &BudgetDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:    reqBody.BudgetsID,
		ActivitiesID: reqBody.ActivitiesID,
	}
	message, err := s.checkBudgetsDetailsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	budgetDetail, err := s.Storage.BudgetDetailsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budgetDetail)

}

func (s *APIServer) UpdateBudgetDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:    reqBody.BudgetsID,
		ActivitiesID: reqBody.ActivitiesID,
	}

	message, err := s.checkBudgetsDetailsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetDetail, err := s.Storage.BudgetDetailsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if updatedBudgetDetail == nil {
		return respondWithError(requestLog, "data budget details not found", err)
	}
	return respondWithSuccess(requestLog, updatedBudgetDetail)

}

func (s *APIServer) DeleteBudgetDetail(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedBudgetDetail, err := s.Storage.BudgetDetailsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if deletedBudgetDetail == nil {
		return respondWithError(requestLog, "data budget details not found", err)
	}

	return respondWithSuccess(requestLog, deletedBudgetDetail)

}
