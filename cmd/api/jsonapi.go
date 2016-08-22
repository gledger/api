package main

type jsonApiDocument struct {
	Data interface{} `json:"data"`
}

type jsonApiTransactionResource struct {
	Type          string                               `json:"type"`
	Id            string                               `json:"id"`
	Attributes    jsonApiTransactionResourceAttributes `json:"attributes"`
	Relationships jsonApiTransactionsRelationships     `json:"relationships"`
}

type jsonApiTransactionResourceAttributes struct {
	Payee        string `json:"payee"`
	Amount       int64  `json:"amount"`
	RollingTotal int64  `json:"rolling_total"`
	OccurredAt   Date   `json:"occurred_at"`
	Cleared      bool   `json:"cleared"`
	Reconciled   bool   `json:"reconciled"`
}

type jsonApiTransactionsRelationships struct {
	Account jsonApiTransactionsRelationshipsAccount `json:"account"`
}

type jsonApiTransactionsRelationshipsAccount struct {
	Data jsonApiAccountResource `json:"data"`
}

type jsonApiAccountResource struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}
