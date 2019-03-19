package app

import (
	"time"

	"github.com/fractalpal/eventflow-example/models/payment"
)

type Payment payment.Payment

type ThirdParty struct {
	payment.ThirdParty
	Time time.Time `json:"-" bson:"-"`
}
