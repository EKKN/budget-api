package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *APIServer) HandlerUserLogin(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	AppLog("username login")
	bodyBytes, requestLog, err := s.prepareRequest(r)
	if err != nil {
		responseLog := LogResponse("error", nil, "failed to prepare request "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to prepare request")
	}

	reqBody := &Users{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		responseLog := LogResponse("error", nil, "failed to decode request body "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("failed to decode request body")
	}

	user, err := s.Storage.UsersStorage.GetDataByLogin(reqBody)
	if err != nil {
		responseLog := LogResponseError("error", "error DB "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error DB")
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	tokenJwt, err := CreateJwt(user)
	if err != nil {
		responseLog := LogResponseError("error", "create jwt "+err.Error())
		AppLog(LogRequestResponse(requestLog, responseLog))
		return nil, fmt.Errorf("error token")
	}

	responseLog := map[string]interface{}{
		"status": "success",
	}

	AppLog(LogRequestResponse(requestLog, responseLog))
	responseLog["token"] = tokenJwt
	return responseLog, nil
}
