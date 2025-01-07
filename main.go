package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	SERVER_PORT := os.Getenv("SERVER_PORT")
	mysql, err := NewMysql()
	if err != nil {
		log.Fatal("Error creating MySQL connection:", err)
	}
	defer mysql.db.Close()

	activitiesStorage := NewActivitiesStorage(mysql)
	usersStorage := NewUsersStorage(mysql)
	budgetPostsStorage := NewBudgetPostsStorage(mysql)
	budgetCapsStorage := NewBudgetCapsStorage(mysql)
	budgetsStorage := NewBudgetsStorage(mysql)
	budgetDetailsStorage := NewBudgetDetailsStorage(mysql)
	budgetDetailsPostsStorage := NewBudgetDetailsPostsStorage(mysql)
	fundRequestsStorage := NewFundRequestsStorage(mysql)
	fundRequestDetailsStorage := NewFundRequestDetailsStorage(mysql)
	budgetDetailsPostsRecommendationsStorage := NewBudgetDetailsPostsRecommendationsStorage(mysql)

	storage := &Storage{
		ActivitiesStorage:                        activitiesStorage,
		UsersStorage:                             usersStorage,
		BudgetPostsStorage:                       budgetPostsStorage,
		BudgetCapsStorage:                        budgetCapsStorage,
		BudgetsStorage:                           budgetsStorage,
		BudgetDetailsStorage:                     budgetDetailsStorage,
		BudgetDetailsPostsStorage:                budgetDetailsPostsStorage,
		FundRequestsStorage:                      fundRequestsStorage,
		FundRequestDetailsStorage:                fundRequestDetailsStorage,
		BudgetDetailsPostsRecommendationsStorage: budgetDetailsPostsRecommendationsStorage,
	}
	AppLog("service run on port ", SERVER_PORT)
	server := NewAPIServer(SERVER_PORT, storage)
	server.Run()

}
