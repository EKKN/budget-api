package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetsRequest(reqBody *Budgets) error {
	if reqBody.Name == "" {
		return fmt.Errorf("name must be filled")
	} else if len(reqBody.Name) > 255 {
		return fmt.Errorf("max length name 255")
	} else if len(reqBody.Description) > 255 {
		return fmt.Errorf("max length description 255")
	} else if reqBody.Periode == "" {
		return fmt.Errorf("periode must be filled")
	} else if reqBody.UnitsID == 0 {
		return fmt.Errorf("unitsID must be filled and valid")
	}
	return nil
}

func (s *APIServer) GetAllBudgets(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgets, err := s.Storage.BudgetsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgets)

}

func (s *APIServer) GetBudgetByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budget, err := s.Storage.BudgetsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budget)

}

func (s *APIServer) CreateBudget(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &Budgets{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	budgetByName, err := s.Storage.BudgetsStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if budgetByName != nil {
		return respondWithError(requestLog, "name already in use", nil)
	}

	budget, err := s.Storage.BudgetsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budget)

}

func (s *APIServer) UpdateBudget(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &Budgets{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	budgetByName, err := s.Storage.BudgetsStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if budgetByName != nil && budgetByName.ID != id {
		return respondWithError(requestLog, "name already in use", nil)
	}

	updatedBudget, err := s.Storage.BudgetsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "error updating budget", err)
	}

	if updatedBudget == nil {
		return respondWithError(requestLog, "data budgets not found", err)
	}
	return respondWithSuccess(requestLog, updatedBudget)

}

func (s *APIServer) DeleteBudget(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedBudget, err := s.Storage.BudgetsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "error deleting budget", err)
	}
	if deletedBudget == nil {
		return respondWithError(requestLog, "data budgets not found", err)
	}
	return respondWithSuccess(requestLog, deletedBudget)

}

func (s *APIServer) UpdateBudgetApproval(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &Budgets{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	updatedBudget, err := s.Storage.BudgetsStorage.UpdateApproved(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "error updating budget approval", err)
	}

	return respondWithSuccess(requestLog, updatedBudget)

}
