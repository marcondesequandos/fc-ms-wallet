package main

import (
	"database/sql"
	"fmt"

	"github.com.br/fc-ms-wallet/internal/database"
	"github.com.br/fc-ms-wallet/internal/event"
	"github.com.br/fc-ms-wallet/internal/usecase/create_account"
	"github.com.br/fc-ms-wallet/internal/usecase/create_client"
	"github.com.br/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com.br/fc-ms-wallet/internal/web"
	"github.com.br/fc-ms-wallet/internal/web/webserver"
	"github.com.br/fc-ms-wallet/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3308", "wallet"))
	fmt.Println("err", err)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)
	// fmt.Println("clientDb", clientDb)
	// fmt.Println("accDb", accountDb)
	// fmt.Println("transDb", transactionDb)

	createClientUseCase := create_client.NewCreateClientUSeCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)
	// fmt.Println("clientuc", createClientUseCase)
	// fmt.Println("accuc", createAccountUseCase)
	// fmt.Println("transuc", createTransactionUseCase)

	webserver := webserver.NewWebServer(":3000")
	fmt.Println("ws", webserver)

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)
	// fmt.Println("clientHandler", clientHandler)
	// fmt.Println("accHandler", accountHandler)
	// fmt.Println("transHandler", transactionHandler)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)
	fmt.Println("Server is running")

	webserver.Start()
}
