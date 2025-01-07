package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIError struct {
	JobId   string `json:"jobid"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func WriteAPISuccess(w http.ResponseWriter, v interface{}, jobID string) {
	data, ok := v.(map[string]interface{})
	if !ok {
		WriteJSON(w, http.StatusBadRequest, APIError{
			Status:  "error",
			JobId:   jobID,
			Message: "Invalid type",
		})

		return
	}

	data["JobId"] = jobID
	WriteJSON(w, http.StatusOK, v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func MakeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := JobID()
		r.Header.Set("jobID", jobID)
		result, err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobId:   jobID,
				Message: err.Error(),
			})
			return
		}
		WriteAPISuccess(w, result, jobID)
	}
}

type APIServer struct {
	ListenAddr string
	Storage    Storage
}

func NewAPIServer(listenNewAddr string, storage *Storage) *APIServer {
	return &APIServer{
		ListenAddr: listenNewAddr,
		Storage:    *storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	publicRouter := router.PathPrefix("/user").Subrouter()
	publicRouter.HandleFunc("/login", MakeHttpHandleFunc(s.HandlerUserLogin)).Methods("POST")

	budgetsRouter := router.PathPrefix("/budgets").Subrouter()
	budgetsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetsGetData)).Methods("GET")
	budgetsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetsCreate)).Methods("POST")
	budgetsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetsGetDataById)).Methods("GET")
	budgetsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetsUpdate)).Methods("PUT")
	budgetsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetsDelete)).Methods("DELETE")
	budgetsRouter.HandleFunc("/approve/{id}", MakeHttpHandleFunc(s.HandlerBudgetsUpdateApproveById)).Methods("PUT")

	activitiesRouter := router.PathPrefix("/activities").Subrouter()
	activitiesRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerActivitiesGetData)).Methods("GET")
	activitiesRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerActivitiesCreate)).Methods("POST")
	activitiesRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerActivitiesGetDataById)).Methods("GET")
	activitiesRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerActivitiesUpdate)).Methods("PUT")
	activitiesRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerActivitiesDelete)).Methods("DELETE")
	activitiesRouter.HandleFunc("/active/{id}", MakeHttpHandleFunc(s.HandlerActivitiesUpdateActiveById)).Methods("PUT")

	budgetPostsRouter := router.PathPrefix("/budget-posts").Subrouter()
	budgetPostsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetPostsGetData)).Methods("GET")
	budgetPostsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetPostsCreate)).Methods("POST")
	budgetPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetPostsGetDataById)).Methods("GET")
	budgetPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetPostsUpdate)).Methods("PUT")
	budgetPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetPostsDelete)).Methods("DELETE")
	budgetPostsRouter.HandleFunc("/active/{id}", MakeHttpHandleFunc(s.HandlerBudgetPostsUpdateActiveById)).Methods("PUT")

	budgetCapsRouter := router.PathPrefix("/budget-caps").Subrouter()
	budgetCapsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetCapsGetData)).Methods("GET")
	budgetCapsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetCapsCreate)).Methods("POST")
	budgetCapsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetCapsGetDataById)).Methods("GET")
	budgetCapsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetCapsUpdate)).Methods("PUT")
	budgetCapsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetCapsDelete)).Methods("DELETE")

	budgetDetailsRouter := router.PathPrefix("/budget-details").Subrouter()
	budgetDetailsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsGetData)).Methods("GET")
	budgetDetailsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsCreate)).Methods("POST")
	budgetDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsGetDataById)).Methods("GET")
	budgetDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsUpdate)).Methods("PUT")
	budgetDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsDelete)).Methods("DELETE")

	budgetDetailsPostsRouter := router.PathPrefix("/budget-details-posts").Subrouter()
	budgetDetailsPostsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsGetData)).Methods("GET")
	budgetDetailsPostsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsCreate)).Methods("POST")
	budgetDetailsPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsGetDataById)).Methods("GET")
	budgetDetailsPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsUpdate)).Methods("PUT")
	budgetDetailsPostsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsDelete)).Methods("DELETE")

	fundRequestsRouter := router.PathPrefix("/fund-requests").Subrouter()
	fundRequestsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerFundRequestsGetData)).Methods("GET")
	fundRequestsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerFundRequestsCreate)).Methods("POST")
	fundRequestsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestsGetDataById)).Methods("GET")
	fundRequestsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestsUpdate)).Methods("PUT")
	fundRequestsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestsDelete)).Methods("DELETE")

	fundRequestDetailsRouter := router.PathPrefix("/fund-request-details").Subrouter()
	fundRequestDetailsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerFundRequestDetailsGetData)).Methods("GET")
	fundRequestDetailsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerFundRequestDetailsCreate)).Methods("POST")
	fundRequestDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestDetailsGetDataById)).Methods("GET")
	fundRequestDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestDetailsUpdate)).Methods("PUT")
	fundRequestDetailsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerFundRequestDetailsDelete)).Methods("DELETE")

	budgetDetailsPostsRecommendationsRouter := router.PathPrefix("/budget-details-posts-recommendations").Subrouter()
	budgetDetailsPostsRecommendationsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsRecommendationsGetData)).Methods("GET")
	budgetDetailsPostsRecommendationsRouter.HandleFunc("", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsRecommendationsCreate)).Methods("POST")
	budgetDetailsPostsRecommendationsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsRecommendationsGetDataById)).Methods("GET")
	budgetDetailsPostsRecommendationsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsRecommendationsUpdate)).Methods("PUT")
	budgetDetailsPostsRecommendationsRouter.HandleFunc("/{id}", MakeHttpHandleFunc(s.HandlerBudgetDetailsPostsRecommendationsDelete)).Methods("DELETE")

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}

func (s *APIServer) SetJobID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) prepareRequest(r *http.Request) ([]byte, map[string]interface{}, error) {
	bodyBytes, err := ReadAndRestoreRequestBody(r)
	if err != nil {
		return nil, nil, err
	}
	requestLog := LogRequest(r, bodyBytes)
	return bodyBytes, requestLog, nil
}

func (s *APIServer) GetID(r *http.Request) (int64, error) {
	strId := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", strId)
	}
	return id, nil
}
