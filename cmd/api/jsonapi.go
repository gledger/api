package main

type jsonAPIDocument struct {
	Data interface{} `json:"data"`
}

type jsonAPITransactionResource struct {
	Type          string                               `json:"type"`
	ID            string                               `json:"id"`
	Attributes    jsonAPITransactionResourceAttributes `json:"attributes"`
	Relationships jsonAPITransactionsRelationships     `json:"relationships"`
}

type jsonAPITransactionResourceAttributes struct {
	Payee        string `json:"payee"`
	Amount       int64  `json:"amount"`
	RollingTotal int64  `json:"rolling_total"`
	OccurredAt   Date   `json:"occurred_at"`
	Cleared      bool   `json:"cleared"`
	Reconciled   bool   `json:"reconciled"`
}

type jsonAPIAccountResource struct {
	Type          string                               `json:"type"`
	ID            string                               `json:"id"`
	Attributes    *jsonAPIAccountResourceAttributes    `json:"attributes,omitempty"`
	Relationships *jsonAPIAccountResourceRelationships `json:"relationships,omitempty"`
}

type jsonAPIAccountResourceAttributes struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Active  bool   `json:"active"`
	Balance int64  `json:"balance"`
}

type jsonAPITransactionsRelationships struct {
	Account  jsonAPITransactionsRelationshipsAccount  `json:"account"`
	Envelope jsonAPITransactionsRelationshipsEnvelope `json:"envelope"`
}

type jsonAPIAccountResourceRelationships struct {
	Transactions map[string]map[string]string `json:"transactions"`
}

type jsonAPITransactionsRelationshipsAccount struct {
	Data jsonAPIAccountResource `json:"data"`
}

type jsonAPITransactionsRelationshipsEnvelope struct {
	Data jsonAPIEnvelopeResource `json:"data"`
}

type jsonAPIEnvelopeResource struct {
	Type       string                             `json:"type"`
	ID         string                             `json:"id"`
	Attributes *jsonAPIEnvelopeResourceAttributes `json:"attributes,omitempty"`
}

type jsonAPIEnvelopeResourceAttributes struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Balance int64  `json:"balance"`
}
