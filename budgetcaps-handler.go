package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetCapsRequest(reqBody *BudgetCaps) error {
	if reqBody.BudgetsID <= 0 {
		return fmt.Errorf("budgets id must be filled")
	} else if reqBody.BudgetPostsID <= 0 {
		return fmt.Errorf("budget posts id must be filled")
	} else if reqBody.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

func (s *APIServer) checkBudgetsCapsForeignKey(primaryKey *PrimaryKeyID) (string, error) {

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:     primaryKey.BudgetsID,
		BudgetPostsID: primaryKey.BudgetPostsID,
	}

	primaryKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKey)
	if err != nil {
		return "error DB", err
	}

	if primaryKey.BudgetsID == 0 {
		return "data budgets not found", fmt.Errorf("data budgets not found")
	}

	if primaryKey.BudgetPostsID == 0 {
		return "data budget post not found", fmt.Errorf("data budget post not found")
	}

	return "ok", nil
}

func (s *APIServer) GetAllBudgetCaps(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgetCaps, err := s.Storage.BudgetCapsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetCaps)

}

func (s *APIServer) GetBudgetCapByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budgetCap, err := s.Storage.BudgetCapsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetCap)

}

func (s *APIServer) CreateBudgetCap(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &BudgetCaps{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetCapsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:     reqBody.BudgetsID,
		BudgetPostsID: reqBody.BudgetPostsID,
	}

	message, err := s.checkBudgetsCapsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	budgetCap, err := s.Storage.BudgetCapsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budgetCap)

}

func (s *APIServer) UpdateBudgetCap(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetCaps{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetCapsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetsID:     reqBody.BudgetsID,
		BudgetPostsID: reqBody.BudgetPostsID,
	}

	message, err := s.checkBudgetsCapsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetCap, err := s.Storage.BudgetCapsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if updatedBudgetCap == nil {
		return respondWithError(requestLog, "data budget cap not found", err)
	}

	return respondWithSuccess(requestLog, updatedBudgetCap)

}

func (s *APIServer) DeleteBudgetCap(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedBudgetCap, err := s.Storage.BudgetCapsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if deletedBudgetCap == nil {
		return respondWithError(requestLog, "data budget cap not found", err)
	}

	return respondWithSuccess(requestLog, deletedBudgetCap)

}
