package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerBudgetDetailsPostsGetData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	budgetDetailsPosts, err := s.Storage.BudgetDetailsPostsStorage.GetData()
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetailsPosts)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsPostsGetDataById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	budgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.GetDataByID(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB %w", err)
	}

	responseLog := LogResponseSuccess(budgetDetailsPost)
	AppLog(LogRequestResponse(requestLog, responseLog))

	return responseLog, nil
}

func validateBudgetDetailsPostsRequest(reqBody *BudgetDetailsPosts) error {

	if reqBody.BudgetDetailsID <= 0 {
		return fmt.Errorf("budget details id must grater than 0")
	} else if reqBody.BudgetPostsID <= 0 {
		return fmt.Errorf("budget post id must grater than 0")
	} else if reqBody.PlannedAmount <= 0 {
		return fmt.Errorf("planned amount must grater than 0")
	} else if reqBody.ApprovedAmount <= 0 {
		return fmt.Errorf("approved amount must grater than 0")
	} else if reqBody.UsageAmount <= 0 {
		return fmt.Errorf("usage amount must grater than 0")
	}

	return nil
}

func (s *APIServer) HandlerBudgetDetailsPostsCreate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponseError("error", "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &BudgetDetailsPosts{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetDetailsPostsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	// Check for related BudgetDetails and BudgetPosts
	budgetDetailsByID, err := s.Storage.BudgetDetailsStorage.GetDataByID(reqBody.BudgetDetailsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetDetailsByID == nil {
		responseLog := LogResponseError("error", "data budget details not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget details not found")
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

	budgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Create(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	responseLog := LogResponseSuccess(budgetDetailsPost)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}

func (s *APIServer) HandlerBudgetDetailsPostsUpdate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	reqBody := &BudgetDetailsPosts{}
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		responseLog := LogResponseError("error", "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	if err := validateBudgetDetailsPostsRequest(reqBody); err != nil {
		responseLog := LogResponseError("error", err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, err
	}

	// Check for related BudgetDetails and BudgetPosts
	budgetDetailsByID, err := s.Storage.BudgetDetailsStorage.GetDataByID(reqBody.BudgetDetailsID)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if budgetDetailsByID == nil {
		responseLog := LogResponseError("error", "data budget details not found")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("data budget details not found")
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

	updatedBudgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Update(id, reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}

	if updatedBudgetDetailsPost == nil {
		responseLog := LogResponseError("error", "no data found to update")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to update")
	}
	responseLog := LogResponseSuccess(updatedBudgetDetailsPost)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
func (s *APIServer) HandlerBudgetDetailsPostsDelete(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	deletedBudgetDetailsPost, err := s.Storage.BudgetDetailsPostsStorage.Delete(id)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if deletedBudgetDetailsPost == nil {
		responseLog := LogResponseError("error", "no data found to delete")
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("no data found to delete")
	}

	responseLog := LogResponseSuccess(deletedBudgetDetailsPost)
	AppLog(LogRequestResponse(requestLog, responseLog))
	return responseLog, nil
}
