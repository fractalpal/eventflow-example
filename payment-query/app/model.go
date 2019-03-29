package app

import (
	"github.com/fractalpal/eventflow-example/models/payment"
)

// Payment model
type Payment payment.Payment

// ThirdParty model
type ThirdParty struct {
	payment.ThirdParty
	Timestamp int64 `json:"-" bson:"-"`
}
