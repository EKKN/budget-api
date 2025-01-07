package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerBudgetDetailsPostsRecommendationsGetData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	budgetDetailsPostsRecommendations, err := s.Storage.BudgetDetailsPostsRecommendationsStorage.GetData()
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetailsPostsRecommendations)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsPostsRecommendationsGetDataById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	id, err := s.GetID(r)
	if err != nil {
		responseLog := LogResponseError("error", "invalid ID "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("invalid ID")
	}
	budgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailsPostsRecommendationsStorage.GetDataId(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB %w", err)
	}

	responseLog := LogResponseSuccess(budgetDetailsPostsRecommendation)
	AppLog(LogRequestResponse(requestLog, responseLog))

	return responseLog, nil
}

func validateBudgetDetailsPostsRecommendationsRequest(reqBody *BudgetDetailsPostsRecommendations) error {
	if reqBody.BudgetDetailsPostsID <= 0 {
		return fmt.Errorf("budget details posts id must be greater than 0")
	} else if reqBody.UserGroupsID <= 0 {
		return fmt.Errorf("user groups id must be greater than 0")
	} else if reqBody.Recommendation <= 0 {
		return fmt.Errorf("recommendation must be greater than 0")
	}
	return nil
}

func (s *APIServer) HandlerBudgetDetailsPostsRecommendationsCreate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &BudgetDetailsPostsRecommendations{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetDetailsPostsRecommendationsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	budgetDetailsPostsById, err := s.Storage.BudgetDetailsPostsStorage.GetDataByID(reqBody.BudgetDetailsPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetDetailsPostsById == nil {
		responseLog := LogResponseError("error", "data budget details posts not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget details posts not found")
	}

	budgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailsPostsRecommendationsStorage.Create(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetailsPostsRecommendation)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsPostsRecommendationsUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	id, err := s.GetID(r)
	if err != nil {
		responseLog := LogResponseError("error", "invalid ID "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("invalid ID")
	}

	reqBody := &BudgetDetailsPostsRecommendations{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetDetailsPostsRecommendationsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	budgetDetailsPostsById, err := s.Storage.BudgetDetailsPostsStorage.GetDataByID(reqBody.BudgetDetailsPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetDetailsPostsById == nil {
		responseLog := LogResponseError("error", "data budget details posts not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget details posts not found")
	}

	updatedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailsPostsRecommendationsStorage.Update(id, reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if updatedBudgetDetailsPostsRecommendation == nil {
		responseLog := LogResponseError("error", "no data found to update")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to update")
	}

	responseLog := LogResponse("success", updatedBudgetDetailsPostsRecommendation, "")
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsPostsRecommendationsDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	id, err := s.GetID(r)
	if err != nil {
		responseLog := LogResponseError("error", "invalid ID "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("invalid ID")
	}

	deletedBudgetDetailsPostsRecommendation, err := s.Storage.BudgetDetailsPostsRecommendationsStorage.Delete(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if deletedBudgetDetailsPostsRecommendation == nil {
		responseLog := LogResponseError("error", "no data found to delete")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to delete")
	}

	responseLog := LogResponseSuccess(deletedBudgetDetailsPostsRecommendation)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
