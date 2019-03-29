package app

import (
	"github.com/fractalpal/eventflow-example/models/payment"
)

type Payment payment.Payment

type ThirdParty struct {
	payment.ThirdParty
	Timestamp int64 `json:"-" bson:"-"`
}
