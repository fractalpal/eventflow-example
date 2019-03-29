package payment

type ThirdParty struct {
	PaymentID     string `json:"-" bson:"-"`
	AccountName   string `json:"account_name" bson:"account_name"`
	AccountNumber string `json:"account_number" bson:"account_number"`
}

type Attributes struct {
	Amount           string     `json:"amount" bson:"amount"`
	Currency         string     `json:"currency" bson:"currency"`
	BeneficiaryParty ThirdParty `json:"beneficiary_party" bson:"beneficiary_party"`
	DebtorParty      ThirdParty `json:"debtor_party" bson:"debtor_party"`
	PaymentID        string     `json:"payment_id" bson:"payment_id"`
	PaymentType      string     `json:"payment_type" bson:"payment_type"`
	ProcessingDate   string     `json:"processing_date" bson:"processing_date"`
	Reference        string     `json:"reference" bson:"reference"`
}

type Payment struct {
	Type                string     `json:"type" bson:"type"`
	ID                  string     `json:"id" bson:"id"`
	Version             int        `json:"version" bson:"version"`
	Attributes          Attributes `json:"attributes" bson:"attributes"`
	LastUpdateTimestamp int64      `json:"last_update_timestamp" bson:"last_update_timestamp"`
}
