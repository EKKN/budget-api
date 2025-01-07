package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerBudgetCapsGetData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	budgetCaps, err := s.Storage.BudgetCapsStorage.GetData()
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetCaps)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetCapsGetDataById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	budgetCap, err := s.Storage.BudgetCapsStorage.GetDataId(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB %w", err)
	}

	responseLog := LogResponseSuccess(budgetCap)
	AppLog(LogRequestResponse(requestLog, responseLog))

	return responseLog, nil
}

func validateBudgetCapsRequest(reqBody *BudgetCaps) error {
	if reqBody.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	return nil

}

func (s *APIServer) HandlerBudgetCapsCreate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &BudgetCaps{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body1 "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetCapsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
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

	budgetPostByID, err := s.Storage.BudgetPostsStorage.GetDataId(reqBody.BudgetPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetPostByID == nil {
		responseLog := LogResponseError("error", "data budget post not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget post not found")
	}

	budgetCap, err := s.Storage.BudgetCapsStorage.Create(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")

	}

	responseLog := LogResponseSuccess(budgetCap)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetCapsUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	reqBody := &BudgetCaps{}
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetCapsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
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

	budgetPostByID, err := s.Storage.BudgetPostsStorage.GetDataId(reqBody.BudgetPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetPostByID == nil {
		responseLog := LogResponseError("error", "data budget post not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget post not found")
	}

	updatedBudgetCap, err := s.Storage.BudgetCapsStorage.Update(id, reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if updatedBudgetCap == nil {
		responseLog := LogResponseError("error", "no data found to update")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to update")
	}
	responseLog := LogResponseSuccess(updatedBudgetCap)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetCapsDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	deletedBudgetCap, err := s.Storage.BudgetCapsStorage.Delete(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if deletedBudgetCap == nil {
		responseLog := LogResponseError("error", "no data found to delete")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to delete")
	}

	responseLog := LogResponseSuccess(deletedBudgetCap)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
