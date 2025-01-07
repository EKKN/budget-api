package main

import "time"

type Users struct {
	ID       int    `json:"id"`
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

type Activities struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description" validate:"required"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Budgets struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Periode     string    `json:"periode"`
	IsApproved  bool      `json:"is_approved"`
	UnitsID     int64     `json:"units_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BudgetPosts struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BudgetCaps struct {
	ID            int64     `json:"id"`
	BudgetsID     int64     `json:"budgets_id"`
	BudgetPostsID int64     `json:"budget_posts_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BudgetDetails struct {
	ID           int64     `json:"id"`
	BudgetsID    int64     `json:"budgets_id"`
	ActivitiesID int64     `json:"activities_id"`
	Description  string    `json:"description"`
	Target       time.Time `json:"target"`
	Quantity     float64   `json:"quantity"`
	UnitValue    float64   `json:"unit_value"`
	Total        float64   `json:"total"`
	Terms        float64   `json:"terms"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BudgetDetailsPosts struct {
	ID              int64     `json:"id"`
	BudgetDetailsID int64     `json:"budget_details_id"`
	BudgetPostsID   int64     `json:"budget_posts_id"`
	PlannedAmount   float64   `json:"planned_amount"`
	ApprovedAmount  float64   `json:"approved_amount"`
	UsageAmount     float64   `json:"usage_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type FundRequests struct {
	ID            int64     `json:"id"`
	BudgetPostsID int64     `json:"budget_posts_id"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FundRequestDetails struct {
	ID              int64     `json:"id"`
	FundRequestsID  int64     `json:"fund_requests_id"`
	ActivitiesID    int64     `json:"activities_id"`
	BudgetDetailsID int64     `json:"budget_details_id"`
	Amount          float64   `json:"amount"`
	Recommendation  string    `json:"recommendation"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BudgetDetailsPostsRecommendations struct {
	ID                   int64     `json:"id"`
	BudgetDetailsPostsID int64     `json:"budget_details_posts_id"`
	UserGroupsID         int64     `json:"user_groups_id"`
	Recommendation       float64   `json:"recommendation"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
