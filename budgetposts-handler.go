package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateBudgetPostsRequest(reqBody *BudgetPosts) error {
	if reqBody.Name == "" {
		return fmt.Errorf("name must be filled")
	} else if len(reqBody.Name) > 255 {
		return fmt.Errorf("max length name 255")
	} else if len(reqBody.Description) > 255 {
		return fmt.Errorf("max length description 255")
	}
	return nil
}

func (s *APIServer) validateBudgetPostsForeignKey(primaryKey *PrimaryKeyID, validateSelfID bool) (string, error) {

	newKey := &PrimaryKeyID{
		BudgetPostsID: primaryKey.BudgetPostsID,
	}

	storedKey, err := s.Storage.PrimaryKeyIDStorage.GetPrimaryKey(newKey)
	if err != nil {
		return "database error", err
	}

	if validateSelfID && storedKey.BudgetPostsID == 0 {
		return "budget posts not found", fmt.Errorf("budget posts not found")
	}

	return "ok", nil
}

func (s *APIServer) GetAllBudgetPosts(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	budgetPosts, err := s.Storage.BudgetPostsStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetPosts)

}

func (s *APIServer) GetBudgetPostByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	budgetPost, err := s.Storage.BudgetPostsStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, budgetPost)

}

func (s *APIServer) CreateBudgetPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &BudgetPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetPostsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	budgetPostByName, err := s.Storage.BudgetPostsStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if budgetPostByName != nil {
		return respondWithError(requestLog, "name already in use", nil)
	}

	budgetPost, err := s.Storage.BudgetPostsStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, budgetPost)

}

func (s *APIServer) UpdateBudgetPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateBudgetPostsRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: id,
	}

	message, err := s.validateBudgetPostsForeignKey(newPrimaryKey, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	budgetPostByName, err := s.Storage.BudgetPostsStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if budgetPostByName != nil && budgetPostByName.ID != id {
		return respondWithError(requestLog, "name already in use", nil)
	}

	updatedBudgetPost, err := s.Storage.BudgetPostsStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	// if updatedBudgetPost == nil {
	// 	return respondWithError(requestLog, "data budget posts not found", err)
	// }
	return respondWithSuccess(requestLog, updatedBudgetPost)

}

func (s *APIServer) DeleteBudgetPost(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: id,
	}

	message, err := s.validateBudgetPostsForeignKey(newPrimaryKey, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	deletedBudgetPost, err := s.Storage.BudgetPostsStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	// if deletedBudgetPost == nil {
	// 	return respondWithError(requestLog, "data budget posts not found", err)
	// }

	return respondWithSuccess(requestLog, deletedBudgetPost)

}

func (s *APIServer) UpdateBudgetPostActiveByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &BudgetPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	newPrimaryKey := &PrimaryKeyID{
		BudgetPostsID: id,
	}

	message, err := s.validateBudgetPostsForeignKey(newPrimaryKey, true)
	if err != nil {
		return respondWithError(requestLog, message, err)
	}

	updatedBudgetPost, err := s.Storage.BudgetPostsStorage.UpdateActive(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, updatedBudgetPost)

}
