package main

import "time"

type Users struct {
	ID       int    `json:"id"`
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

type Activities struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
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

// type BudgetCaps struct {
// 	ID            int       `json:"id"`
// 	BudgetsID     string    `json:"budgets_id"`
// 	BudgetPostsID string    `json:"budget_posts_id"`
// 	Amount        float64   `json:"amount"`
// 	CreatedAt     time.Time `json:"created_at"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// }

// BudgetCaps represents the budget_caps table structure.
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
	UsageAmount     string    `json:"usage_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
