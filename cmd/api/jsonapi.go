package main

type jsonApiDocument struct {
	Data interface{} `json:"data"`
}

type jsonApiTransactionResource struct {
	Type       string                               `json:"type"`
	Id         string                               `json:"id"`
	Attributes jsonApiTransactionResourceAttributes `json:"attributes"`
}

type jsonApiTransactionResourceAttributes struct {
	Payee      string `json:"payee"`
	Amount     int64  `json:"amount"`
	OccurredAt Date   `json:"occurred_at"`
	Cleared    bool   `json:"cleared"`
	Reconciled bool   `json:"reconciled"`
}
