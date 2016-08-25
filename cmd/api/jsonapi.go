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

type jsonApiAccountResource struct {
	Type          string                               `json:"type"`
	Id            string                               `json:"id"`
	Attributes    *jsonApiAccountResourceAttributes    `json:"attributes,omitempty"`
	Relationships *jsonApiAccountResourceRelationships `json:"relationships,omitempty"`
}

type jsonApiAccountResourceAttributes struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Active  bool   `json:"active"`
	Balance int64  `json:"balance"`
}

type jsonApiTransactionsRelationships struct {
	Account jsonApiTransactionsRelationshipsAccount `json:"account"`
}

type jsonApiAccountResourceRelationships struct {
	Transactions map[string]map[string]string `json:"transactions"`
}

type jsonApiTransactionsRelationshipsAccount struct {
	Data jsonApiAccountResource `json:"data"`
}
