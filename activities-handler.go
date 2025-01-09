package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func validateActivityRequest(reqBody *Activities) error {
	if reqBody.Name == "" {
		return fmt.Errorf("name must be filled")
	} else if len(reqBody.Name) > 255 {
		return fmt.Errorf("max length name 255")
	} else if len(reqBody.Description) > 255 {
		return fmt.Errorf("max length description 255")
	}
	return nil
}

func (s *APIServer) GetAllActivities(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	activities, err := s.Storage.ActivitiesStorage.GetAll()
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, activities)

}

func (s *APIServer) GetActivityByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	activity, err := s.Storage.ActivitiesStorage.GetById(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	return respondWithSuccess(requestLog, activity)

}

func (s *APIServer) CreateActivity(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	reqBody := &Activities{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateActivityRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	activityByName, err := s.Storage.ActivitiesStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if activityByName != nil {
		return respondWithError(requestLog, "name already in use", nil)
	}

	activity, err := s.Storage.ActivitiesStorage.Create(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	return respondWithSuccess(requestLog, activity)

}

func (s *APIServer) UpdateActivity(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &Activities{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	if err := validateActivityRequest(reqBody); err != nil {
		return respondWithError(requestLog, err.Error(), nil)
	}

	activityByName, err := s.Storage.ActivitiesStorage.GetByName(reqBody.Name)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if activityByName != nil && activityByName.ID != id {
		return respondWithError(requestLog, "name already in use", nil)
	}

	updatedActivity, err := s.Storage.ActivitiesStorage.Update(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if updatedActivity == nil {
		return respondWithError(requestLog, "data activities not found", err)
	}

	return respondWithSuccess(requestLog, updatedActivity)

}

func (s *APIServer) DeleteActivity(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	deletedActivity, err := s.Storage.ActivitiesStorage.Delete(id)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if deletedActivity == nil {
		return respondWithError(requestLog, "data activities not found", err)
	}
	return respondWithSuccess(requestLog, deletedActivity)

}

func (s *APIServer) UpdateActivityStatusByID(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	id, err := s.GetID(r)
	if err != nil {
		return respondWithError(requestLog, "invalid ID", err)
	}

	reqBody := &Activities{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "invalid data request", err)
	}

	updatedActivity, err := s.Storage.ActivitiesStorage.UpdateActive(id, reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}

	if updatedActivity == nil {
		return respondWithError(requestLog, "data activities not found", err)
	}
	return respondWithSuccess(requestLog, updatedActivity)

}
