package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerFundRequestsGetData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	fundRequest, err := s.Storage.FundRequestsStorage.GetData()
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(fundRequest)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerFundRequestsGetDataById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	fundRequest, err := s.Storage.FundRequestsStorage.GetDataId(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB %w", err)
	}

	responseLog := LogResponseSuccess(fundRequest)
	AppLog(LogRequestResponse(requestLog, responseLog))

	return responseLog, nil
}

func validateFundRequestsRequest(reqBody *FundRequests) error {
	if reqBody.BudgetPostsID <= 0 {
		return fmt.Errorf("budget posts id must be greater than 0")
	} else if reqBody.Date.IsZero() {
		return fmt.Errorf("date must be filled")
	} else if reqBody.Type == "" {
		return fmt.Errorf("type must be filled")
	} else if reqBody.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	} else if reqBody.Status <= "" {
		return fmt.Errorf("status must be filled")
	}
	return nil
}

func (s *APIServer) HandlerFundRequestsCreate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &FundRequests{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateFundRequestsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	budgetPostsByid, err := s.Storage.BudgetPostsStorage.GetDataId(reqBody.BudgetPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetPostsByid == nil {
		responseLog := LogResponseError("error", "data budget post not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget post not found")
	}

	fundRequest, err := s.Storage.FundRequestsStorage.Create(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(fundRequest)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerFundRequestsUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	reqBody := &FundRequests{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateFundRequestsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	budgetPostsByid, err := s.Storage.BudgetPostsStorage.GetDataId(reqBody.BudgetPostsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetPostsByid == nil {
		responseLog := LogResponseError("error", "data budget post not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget post not found")
	}

	updatedFundRequest, err := s.Storage.FundRequestsStorage.Update(id, reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if updatedFundRequest == nil {
		responseLog := LogResponseError("error", "no data found to update")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to update")
	}

	responseLog := LogResponse("success", updatedFundRequest, "")
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerFundRequestsDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	deletedFundRequest, err := s.Storage.FundRequestsStorage.Delete(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if deletedFundRequest == nil {
		responseLog := LogResponseError("error", "no data found to delete")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to delete")
	}

	responseLog := LogResponseSuccess(deletedFundRequest)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
