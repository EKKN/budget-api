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

func (s *APIServer) validateBDPRFForeignKey(primaryKey *PrimaryKeyID, validateSelfID bool, checkSelfOnly bool) (string, error) {

	newKey := &PrimaryKeyID{
		BudgetDetailsPostsRecommendationsID: primaryKey.BudgetDetailsPostsRecommendationsID,
		BudgetDetailsPostsID:                primaryKey.BudgetDetailsPostsID,
	}

	storedKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newKey)
	if err != nil {
		return "database error", err
	}

	if validateSelfID && storedKey.BudgetDetailsPostsRecommendationsID == 0 {
		return "budget detail post recommendation not found", fmt.Errorf("budget detail post recommendation not found")
	}
	if !checkSelfOnly {

		if storedKey.BudgetDetailsPostsID == 0 {
			return "data budget details posts not found", fmt.Errorf("data budget details posts not found")
		}
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

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsPostsID: reqBody.BudgetDetailsPostsID,
	}

	message, err := s.validateBDPRFForeignKey(newPrimaryKey, false, false)
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

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsPostsRecommendationsID: id,
		BudgetDetailsPostsID:                reqBody.BudgetDetailsPostsID,
	}

	message, err := s.validateBDPRFForeignKey(newPrimaryKey, true, false)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	// if updatedBudgetDetailsPostsRecommendation == nil {
	// 	return respondWithError(requestLog, "data budget details posts recommendation not found", err)
	// }

	return respondWithSuccess(requestLog, updatedBudgetDetailsPostsRecommendation)

}

func (s *APIServer) DeleteBudgetDetailPostRec(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetDetailsPostsRecommendationsID: id,
	}

	message, err := s.validateBDPRFForeignKey(newPrimaryKey, true, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	deletedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailPostRecStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	// if deletedBudgetDetailsPostsRecommendation == nil {
	// 	return respondWithError(requestLog, "data budget details posts recommendation not found", err)
	// }

	return respondWithSuccess(requestLog, deletedBudgetDetailsPostsRecommendation)

}
