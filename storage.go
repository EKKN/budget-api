package main

type Storage struct {
	UsersStorage                             UsersStorage
	BudgetsStorage                           BudgetsStorage
	ActivitiesStorage                        ActivitiesStorage
	BudgetPostsStorage                       BudgetPostsStorage
	BudgetCapsStorage                        BudgetCapsStorage
	BudgetDetailsStorage                     BudgetDetailsStorage
	BudgetDetailsPostsStorage                BudgetDetailsPostsStorage
	FundRequestsStorage                      FundRequestsStorage
	FundRequestDetailsStorage                FundRequestDetailsStorage
	BudgetDetailsPostsRecommendationsStorage BudgetDetailsPostsRecommendationsStorage
}
