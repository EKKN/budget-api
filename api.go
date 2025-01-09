package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIError struct {
	JobID   string `json:"jobId"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func WriteAPISuccess(w http.ResponseWriter, data interface{}, jobID string) {
	response, ok := data.(map[string]interface{})
	if !ok {
		WriteJSON(w, http.StatusBadRequest, APIError{
			Status:  "error",
			JobID:   jobID,
			Message: "Invalid data type",
		})
		return
	}
	response["JobID"] = jobID
	WriteJSON(w, http.StatusOK, response)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type APIServer struct {
	ListenAddr string
	Storage    Storage
}

func NewAPIServer(listenAddr string, storage *Storage) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Storage:    *storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.Use(s.addJobid)
	// User routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", s.prepareAndHandleRequest(s.UserLogin)).Methods("POST")

	// Budgets routes
	budgetsRouter := router.PathPrefix("/budgets").Subrouter()
	budgetsRouter.Use(s.Authenticate)
	budgetsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgets)).Methods("GET")
	budgetsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudget)).Methods("POST")
	budgetsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetByID)).Methods("GET")
	budgetsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudget)).Methods("PUT")
	budgetsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudget)).Methods("DELETE")
	budgetsRouter.HandleFunc("/approve/{id}", s.prepareAndHandleRequest(s.UpdateBudgetApproval)).Methods("PUT")

	// Activities routes
	activitiesRouter := router.PathPrefix("/activities").Subrouter()
	activitiesRouter.Use(s.Authenticate)
	activitiesRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllActivities)).Methods("GET")
	activitiesRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateActivity)).Methods("POST")
	activitiesRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetActivityByID)).Methods("GET")
	activitiesRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateActivity)).Methods("PUT")
	activitiesRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteActivity)).Methods("DELETE")
	activitiesRouter.HandleFunc("/active/{id}", s.prepareAndHandleRequest(s.UpdateActivityStatusByID)).Methods("PUT")

	// Budget posts routes
	budgetPostsRouter := router.PathPrefix("/budget-posts").Subrouter()
	budgetPostsRouter.Use(s.Authenticate)
	budgetPostsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgetPosts)).Methods("GET")
	budgetPostsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudgetPost)).Methods("POST")
	budgetPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetPostByID)).Methods("GET")
	budgetPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudgetPost)).Methods("PUT")
	budgetPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudgetPost)).Methods("DELETE")
	budgetPostsRouter.HandleFunc("/active/{id}", s.prepareAndHandleRequest(s.UpdateBudgetPostActiveByID)).Methods("PUT")

	// Budget caps routes
	budgetCapsRouter := router.PathPrefix("/budget-caps").Subrouter()
	budgetCapsRouter.Use(s.Authenticate)
	budgetCapsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgetCaps)).Methods("GET")
	budgetCapsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudgetCap)).Methods("POST")
	budgetCapsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetCapByID)).Methods("GET")
	budgetCapsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudgetCap)).Methods("PUT")
	budgetCapsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudgetCap)).Methods("DELETE")

	// Budget details routes
	budgetDetailsRouter := router.PathPrefix("/budget-details").Subrouter()
	budgetDetailsRouter.Use(s.Authenticate)
	budgetDetailsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgetDetails)).Methods("GET")
	budgetDetailsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudgetDetail)).Methods("POST")
	budgetDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetDetailByID)).Methods("GET")
	budgetDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudgetDetail)).Methods("PUT")
	budgetDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudgetDetail)).Methods("DELETE")

	// Budget details posts routes
	budgetDetailsPostsRouter := router.PathPrefix("/budget-details-posts").Subrouter()
	budgetDetailsPostsRouter.Use(s.Authenticate)
	budgetDetailsPostsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgetDetailPosts)).Methods("GET")
	budgetDetailsPostsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudgetDetailPost)).Methods("POST")
	budgetDetailsPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetDetailPostByID)).Methods("GET")
	budgetDetailsPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudgetDetailPost)).Methods("PUT")
	budgetDetailsPostsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudgetDetailPost)).Methods("DELETE")

	// Fund requests routes
	fundRequestsRouter := router.PathPrefix("/fund-requests").Subrouter()
	fundRequestsRouter.Use(s.Authenticate)
	fundRequestsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllFundRequests)).Methods("GET")
	fundRequestsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateFundRequest)).Methods("POST")
	fundRequestsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetFundRequestByID)).Methods("GET")
	fundRequestsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateFundRequest)).Methods("PUT")
	fundRequestsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteFundRequest)).Methods("DELETE")

	// Fund request details routes
	fundRequestDetailsRouter := router.PathPrefix("/fund-request-details").Subrouter()
	fundRequestDetailsRouter.Use(s.Authenticate)
	fundRequestDetailsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllFundRequestDetails)).Methods("GET")
	fundRequestDetailsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateFundRequestDetail)).Methods("POST")
	fundRequestDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetFundRequestDetailByID)).Methods("GET")
	fundRequestDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateFundRequestDetail)).Methods("PUT")
	fundRequestDetailsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteFundRequestDetail)).Methods("DELETE")

	// Budget details posts recommendations routes
	budgetDetailsPostsRecsRouter := router.PathPrefix("/budget-details-posts-recommendations").Subrouter()
	budgetDetailsPostsRecsRouter.Use(s.Authenticate)
	budgetDetailsPostsRecsRouter.HandleFunc("", s.prepareAndHandleRequest(s.GetAllBudgetDetailPostRecs)).Methods("GET")
	budgetDetailsPostsRecsRouter.HandleFunc("", s.prepareAndHandleRequest(s.CreateBudgetDetailPostRec)).Methods("POST")
	budgetDetailsPostsRecsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.GetBudgetDetailPostRecByID)).Methods("GET")
	budgetDetailsPostsRecsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.UpdateBudgetDetailPostRec)).Methods("PUT")
	budgetDetailsPostsRecsRouter.HandleFunc("/{id}", s.prepareAndHandleRequest(s.DeleteBudgetDetailPostRec)).Methods("DELETE")

	// Handle not found and method not allowed routes
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := r.Header.Get("jobID")
		WriteJSON(w, http.StatusBadRequest, APIError{
			Status:  "error",
			JobID:   jobID,
			Message: "Method Not Allowed",
		})

	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := r.Header.Get("jobID")
		WriteJSON(w, http.StatusBadRequest, APIError{
			Status:  "error",
			JobID:   jobID,
			Message: "Page Not found",
		})

	})

	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}

func (s *APIServer) addJobid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID := JobID()
		r.Header.Set("JobID", jobID)
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) prepareAndHandleRequest(handlerFunc func(http.ResponseWriter, *http.Request, []byte, map[string]interface{}) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := r.Header.Get("jobID")
		bodyBytes, requestLog, err := s.prepareRequest(r)
		if err != nil {
			AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": err.Error()}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: err.Error(),
			})
			return
		}

		data, err := handlerFunc(w, r, bodyBytes, requestLog)
		if err != nil {
			// AppLog(LogRequestResponse(requestLog, map[string]interface{}{"status": "error", "message": err.Error()}))
			WriteJSON(w, http.StatusBadRequest, APIError{
				Status:  "error",
				JobID:   jobID,
				Message: err.Error(),
			})
			return
		}

		WriteAPISuccess(w, data, jobID)
	}
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
	strID := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return id, fmt.Errorf("invalid ID : %s", strID)
	}
	return id, nil
}

func ReadAndRestoreRequestBody(r *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}
