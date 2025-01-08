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

	activitiesStorage := NewActivitiesStorage(mysql.db)
	usersStorage := NewUsersStorage(mysql.db)
	budgetPostsStorage := NewBudgetPostsStorage(mysql.db)
	budgetCapsStorage := NewBudgetCapsStorage(mysql.db)
	budgetsStorage := NewBudgetsStorage(mysql.db)
	budgetDetailsStorage := NewBudgetDetailsStorage(mysql.db)
	budgetDetailsPostsStorage := NewBudgetDetailsPostsStorage(mysql.db)
	fundRequestsStorage := NewFundRequestsStorage(mysql.db)
	fundRequestDetailsStorage := NewFundRequestDetailsStorage(mysql.db)
	budgetDetailPostRecStorage := NewBudgetDetailPostRecStorage(mysql.db)

	primaryKeyIDStorage := NewPrimaryKeyIDStorage(mysql.db)

	storage := &Storage{
		ActivitiesStorage:          activitiesStorage,
		UsersStorage:               usersStorage,
		BudgetPostsStorage:         budgetPostsStorage,
		BudgetCapsStorage:          budgetCapsStorage,
		BudgetsStorage:             budgetsStorage,
		BudgetDetailsStorage:       budgetDetailsStorage,
		BudgetDetailsPostsStorage:  budgetDetailsPostsStorage,
		FundRequestsStorage:        fundRequestsStorage,
		FundRequestDetailsStorage:  fundRequestDetailsStorage,
		BudgetDetailPostRecStorage: budgetDetailPostRecStorage,
		PrimaryKeyIDStorage:        primaryKeyIDStorage,
	}
	AppLog("service run on port ", SERVER_PORT)
	server := NewAPIServer(SERVER_PORT, storage)
	server.Run()
}
