package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetDetailsPostsRecommendationsRequest(reqBody *BudgetDetailsPostsRecommendations) error {
	if reqBody.BudgetDetailsPostsID <= 0 {
		return fmt.Errorf("budget details posts id must be filled")
	} else if reqBody.UserGroupsID <= 0 {
		return fmt.Errorf("user groups id must be greater than 0")
	} else if reqBody.Recommendation <= 0 {
		return fmt.Errorf("recommendation must be greater than 0")
	}
	return nil
}

func (s *APIServer) checkBDPRForeignKey(primaryKey *PrimaryKeyID) (string, error) {

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsPostsID: primaryKey.BudgetDetailsPostsID,
	}

	primaryKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newPrimaryKey)
	if err != nil {
		return "error DB", err
	}

	if primaryKey.BudgetDetailsPostsID == 0 {
		return "data budget details posts not found", fmt.Errorf("data budget details posts not found")
	}

	return "ok", nil
}

func (s *APIServer) GetAllBudgetDetailPostRecs(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgetDetailsPostsRecommendations, err := s.Storage.BudgetDetailPostRecStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetailsPostsRecommendations)

}

func (s *APIServer) GetBudgetDetailPostRecByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetDetailsPostsRecommendation)

}

func (s *APIServer) CreateBudgetDetailPostRec(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &BudgetDetailsPostsRecommendations{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsPostsRecommendationsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{}
	newPrimaryKey.BudgetDetailsPostsID = reqBody.BudgetDetailsPostsID

	message, err := s.checkBDPRForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	budgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budgetDetailsPostsRecommendation)

}

func (s *APIServer) UpdateBudgetDetailPostRec(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetDetailsPostsRecommendations{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetDetailsPostsRecommendationsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{}
	newPrimaryKey.BudgetDetailsPostsID = id

	message, err := s.checkBDPRForeignKey(newPrimaryKey)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if updatedBudgetDetailsPostsRecommendation == nil {
		return respondWithError(requestLog, "data budget details posts recommendation not found", err)
	}

	return respondWithSuccess(requestLog, updatedBudgetDetailsPostsRecommendation)

}

func (s *APIServer) DeleteBudgetDetailPostRec(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if deletedBudgetDetailsPostsRecommendation == nil {
		return respondWithError(requestLog, "data budget details posts recommendation not found", err)
	}

	return respondWithSuccess(requestLog, deletedBudgetDetailsPostsRecommendation)

}
