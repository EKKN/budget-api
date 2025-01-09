package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetDetailsPost(reqBody *BudgetDetailsPosts) error {
	if reqBody.BudgetDetailsID <= 0 {
		return fmt.Errorf("budget details id must be filled")
	} else if reqBody.BudgetPostsID <= 0 {
		return fmt.Errorf("budget posts id must be filled")
	} else if reqBody.PlannedAmount <= 0 {
		return fmt.Errorf("planned amount must be greater than 0")
	} else if reqBody.ApprovedAmount <= 0 {
		return fmt.Errorf("approved amount must be greater than 0")
	} else if reqBody.UsageAmount <= 0 {
		return fmt.Errorf("usage amount must be greater than 0")
	}

	return nil
}

func (s *APIServer) checkBudgetsDetailsPostsForeignKey(primaryKey *PrimaryKeyID) (string, error) {

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsID: primaryKey.BudgetDetailsID,
		BudgetPostsID:   primaryKey.BudgetPostsID,
	}

	primaryKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKey)
	if err != nil {
		return "error DB", err
	}

	if primaryKey.BudgetDetailsID == 0 {
		return "data budget details not found", fmt.Errorf("data budget details not found")
	}

	if primaryKey.BudgetPostsID == 0 {
		return "data budget posts found", fmt.Errorf("data budget posts not found")
	}
	return "ok", nil
}

func (s *APIServer) GetAllBudgetDetailPosts(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgetDetailsPosts, err := s.Storage.BudgetDetailsPostsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetailsPosts)

}

func (s *APIServer) GetBudgetDetailPostByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetailsPost)

}

func (s *APIServer) CreateBudgetDetailPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &BudgetDetailsPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsPost(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsID: reqBody.BudgetDetailsID,
		BudgetPostsID:   reqBody.BudgetPostsID,
	}

	message, err := s.checkBudgetsDetailsPostsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	budgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budgetDetailsPost)

}

func (s *APIServer) UpdateBudgetDetailPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetDetailsPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsPost(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsID: reqBody.BudgetDetailsID,
		BudgetPostsID:   reqBody.BudgetPostsID,
	}

	message, err := s.checkBudgetsDetailsPostsForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if updatedBudgetDetailsPost == nil {
		return respondWithError(requestLog, "data budget details not found", err)
	}

	return respondWithSuccess(requestLog, updatedBudgetDetailsPost)

}

func (s *APIServer) DeleteBudgetDetailPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedBudgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if deletedBudgetDetailsPost == nil {
		return respondWithError(requestLog, "data budget details not found", err)
	}

	return respondWithSuccess(requestLog, deletedBudgetDetailsPost)

}
