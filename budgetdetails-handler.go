package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerBudgetDetailsGetData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	budgetDetails, err := s.Storage.BudgetDetailsStorage.GetData()
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetails)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsGetDataById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	budgetDetail, err := s.Storage.BudgetDetailsStorage.GetDataByID(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB %w", err)
	}

	responseLog := LogResponseSuccess(budgetDetail)
	AppLog(LogRequestResponse(requestLog, responseLog))

	return responseLog, nil
}

func validateBudgetDetailsRequest(reqBody *BudgetDetails) error {
	if reqBody.Description == "" {
		return fmt.Errorf("description must be filled")
	}
	return nil
}

func (s *APIServer) HandlerBudgetDetailsCreate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &BudgetDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetDetailsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	ActivitiesByID, err := s.Storage.ActivitiesStorage.GetDataId(reqBody.ActivitiesID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if ActivitiesByID == nil {
		responseLog := LogResponseError("error", "data activities not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data activities not found")
	}

	budgetByID, err := s.Storage.BudgetsStorage.GetDataId(reqBody.BudgetsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetByID == nil {
		responseLog := LogResponseError("error", "data budget not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budgets not found")
	}

	budgetDetail, err := s.Storage.BudgetDetailsStorage.Create(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetail)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	reqBody := &BudgetDetails{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}
	if err := validateBudgetDetailsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	ActivitiesByID, err := s.Storage.ActivitiesStorage.GetDataId(reqBody.ActivitiesID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if ActivitiesByID == nil {
		responseLog := LogResponseError("error", "data activities not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data activities not found")
	}

	budgetByID, err := s.Storage.BudgetsStorage.GetDataId(reqBody.BudgetsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetByID == nil {
		responseLog := LogResponseError("error", "data budget not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budgets not found")
	}

	updatedBudgetDetail, err := s.Storage.BudgetDetailsStorage.Update(id, reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if updatedBudgetDetail == nil {
		responseLog := LogResponseError("error", "no data found to update")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to update")
	}

	responseLog := LogResponse("success", updatedBudgetDetail, "")
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	deletedBudgetDetail, err := s.Storage.BudgetDetailsStorage.Delete(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if deletedBudgetDetail == nil {
		responseLog := LogResponseError("error", "no data found to delete")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to delete")
	}

	responseLog := LogResponseSuccess(deletedBudgetDetail)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
