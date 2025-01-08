package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (s *APIServer) UserLogin(w http.ResponseWriter, r *http.Request, bodyBytes []byte, requestLog map[string]interface{}) (interface{}, error) {

	AppLog("username login")

	reqBody := &Users{}
	if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(reqBody); err != nil {
		return respondWithError(requestLog, "failed to decode request body", err)
	}

	user, err := s.Storage.UsersStorage.GetByLogin(reqBody)
	if err != nil {
		return respondWithError(requestLog, "database error", err)
	}
	if user == nil {
		return respondWithError(requestLog, "user not found", nil)
	}

	tokenJwt, err := CreateJwt(user)
	if err != nil {
		return respondWithError(requestLog, "error creating JWT", err)
	}

	responseLog := map[string]interface{}{
		"status": "success",
		"token":  tokenJwt,
	}

	return respondWithSuccessStruct(requestLog, responseLog)
}
