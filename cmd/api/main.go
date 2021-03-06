package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/net/context"

	"github.com/gledger/api"
	"github.com/gledger/api/db"
)

func main() {
	port := os.Getenv("PORT")
	pgURI := os.Getenv("DATABASE_URL")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Connecting to %s\n", pgURI)
	pg, err := sql.Open("postgres", pgURI)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	ctx := context.Background()
	svc := gledger.NewAccountService(
		db.SaveAccount(pg.Exec),
		db.AllAccounts(pg.Query),
		db.ReadAccount(pg.QueryRow),
	)
	txnSvc := gledger.NewTransactionService(
		db.SaveTransaction(pg),
		db.TransactionsForAccount(pg),
	)
	envSvc := gledger.NewEnvelopeService(
		db.SaveEnvelope(pg.Exec),
		db.AllEnvelopes(pg.Query),
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
		"/transactions",
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

	router.HandleFunc(
		"/envelopes",
		httptransport.NewServer(
			ctx,
			makeCreateEnvelopeEndpoint(envSvc),
			decodeCreateEnvelopeRequest,
			created(encodeResponse),
			httpOptions...,
		).ServeHTTP,
	).Methods("POST")
	router.HandleFunc(
		"/envelopes",
		httptransport.NewServer(
			ctx,
			makeAllEnvelopesEndpoint(envSvc),
			emptyRequest,
			encodeResponse,
			httpOptions...,
		).ServeHTTP,
	).Methods("GET")

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200", "https://gledger-web.herokuapp.com"},
	}).Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		router.ServeHTTP(res, req)
		end := time.Now()
		log.Printf("%s %s %s\n", req.Method, req.URL, end.Sub(start))
	}))
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func emptyRequest(context.Context, *http.Request) (interface{}, error) {
	return nil, nil
}
