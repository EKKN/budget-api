package main

type Storage struct {
	ActivitiesStorage          ActivitiesStorage
	UsersStorage               UsersStorage
	BudgetPostsStorage         BudgetPostsStorage
	BudgetCapsStorage          BudgetCapsStorage
	BudgetsStorage             BudgetsStorage
	BudgetDetailsStorage       BudgetDetailsStorage
	BudgetDetailsPostsStorage  BudgetDetailsPostsStorage
	FundRequestsStorage        FundRequestsStorage
	FundRequestDetailsStorage  FundRequestDetailsStorage
	BudgetDetailPostRecStorage BudgetDetailPostRecStorage
	PrimaryKeyIDStorage        PrimaryKeyIDStorage
}
