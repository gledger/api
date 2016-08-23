package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/net/context"

	"github.com/zombor/gledger"
	"github.com/zombor/gledger/db"
)

func main() {
	port := os.Getenv("PORT")
	pgUri := os.Getenv("DATABASE_URL")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Connecting to %s\n", pgUri)
	pg, err := sql.Open("postgres", pgUri)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	ctx := context.Background()
	svc := gledger.NewAccountService(
		db.SaveAccount(pg),
		db.AllAccounts(pg),
		db.ReadAccount(pg),
	)
	txnSvc := gledger.NewTransactionService(
		db.SaveTransaction(pg),
		db.TransactionsForAccount(pg),
	)

	httpOptions := []httptransport.ServerOption{httptransport.ServerErrorEncoder(errorEncoder)}

	router.HandleFunc(
		"/accounts",
		httptransport.NewServer(
			ctx,
			makeCreateAccountEndpoint(svc),
			decodeCreateAccountRequest,
			created(encodeResponse),
			httpOptions...,
		).ServeHTTP,
	).Methods("POST")
	router.HandleFunc(
		"/accounts",
		httptransport.NewServer(
			ctx,
			makeAllAccountsEndpoint(svc),
			emptyRequest,
			encodeResponse,
			httpOptions...,
		).ServeHTTP,
	).Methods("GET")

	router.HandleFunc(
		"/accounts/{uuid}",
		httptransport.NewServer(
			ctx,
			makeReadAccountEndpoint(svc),
			decodeReadAccountRequest,
			encodeResponse,
			httpOptions...,
		).ServeHTTP,
	).Methods("GET")
	router.HandleFunc(
		"/accounts/{uuid}/transactions",
		httptransport.NewServer(
			ctx,
			makeCreateTransactionEndpoint(txnSvc),
			decodeCreateTransactionRequest,
			created(encodeResponse),
			httpOptions...,
		).ServeHTTP,
	).Methods("POST")

	router.HandleFunc(
		"/accounts/{uuid}/transactions",
		httptransport.NewServer(
			ctx,
			makeReadAccountTransactionsEndpoint(txnSvc),
			decodeReadAccountTransactionsRequest,
			encodeResponse,
			httpOptions...,
		).ServeHTTP,
	).Methods("GET")

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200", "https://gledger-web.herokuapp.com"},
	}).Handler(router)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func emptyRequest(context.Context, *http.Request) (interface{}, error) {
	return nil, nil
}
