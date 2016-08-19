package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
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

	ctx := context.Background()
	svc := gledger.NewAccountService(saveAccount(db), allAccounts(db))

	createAccountHandler := httptransport.NewServer(
		ctx,
		makeCreateAccountEndpoint(svc),
		decodeCreateAccountRequest,
		encodeResponse,
	)
	allAccountsHandler := httptransport.NewServer(
		ctx,
		makeAllAccountsEndpoint(svc),
		func(context.Context, *http.Request) (interface{}, error) {
			return nil, nil
		},
		encodeResponse,
	)

	http.Handle("/create-account", createAccountHandler)
	http.Handle("/all-accounts", allAccountsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
