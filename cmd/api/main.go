package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
)

func main() {
	var pgUri string
	flag.StringVar(&pgUri, "databaseUri", "", "Database URI. Required")
	flag.Parse()

	fmt.Printf("Connecting to %s", pgUri)
	db, err := sql.Open("postgres", pgUri)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	ctx := context.Background()
	svc := gledger.NewAccountService(saveAccount(db), allAccounts(db))
	txnSvc := gledger.NewTransactionService(
		saveTransaction(db),
		transactionsForAccount(db),
	)

	router.HandleFunc(
		"/accounts",
		httptransport.NewServer(
			ctx,
			makeCreateAccountEndpoint(svc),
			decodeCreateAccountRequest,
			encodeResponse,
		).ServeHTTP,
	).Methods("POST")
	router.HandleFunc(
		"/accounts",
		httptransport.NewServer(
			ctx,
			makeAllAccountsEndpoint(svc),
			emptyRequest,
			encodeResponse,
		).ServeHTTP,
	).Methods("GET")

	router.HandleFunc(
		"/accounts/{uuid}/transactions",
		httptransport.NewServer(
			ctx,
			makeCreateTransactionEndpoint(txnSvc),
			decodeCreateTransactionRequest,
			encodeResponse,
		).ServeHTTP,
	).Methods("POST")

	router.HandleFunc(
		"/accounts/{uuid}/transactions",
		httptransport.NewServer(
			ctx,
			makeReadAccountTransactionsEndpoint(txnSvc),
			decodeReadAccountTransactionsRequest,
			encodeResponse,
		).ServeHTTP,
	).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func emptyRequest(context.Context, *http.Request) (interface{}, error) {
	return nil, nil
}
