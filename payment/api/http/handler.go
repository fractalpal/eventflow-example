package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment/app"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	service app.Service
	context context.Context
	logger  logrus.FieldLogger
}

func NewHandler(service app.Service, logger logrus.FieldLogger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fields := logrus.Fields{}
	fields["method"] = []string{"delete"}
	fields["url"] = []string{r.URL.String()}
	l := h.logger.WithFields(fields)
	ctx := log.ContextWithLogger(r.Context(), l)
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(ctx, id); err != nil && err != app.ErrNoRows {
		ctxLogger := log.FromContext(ctx)
		ctxLogger.WithError(err).Error("couldn't delete payment")
		if errors.Cause(err) == app.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l.Debug("payment deleted")

	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	fields := logrus.Fields{}
	fields["method"] = []string{"post"}
	fields["url"] = []string{r.URL.String()}
	l := h.logger.WithFields(fields)
	ctx := log.ContextWithLogger(r.Context(), l)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "couldn't read request body.")
		l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payment app.Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		err = errors.Wrap(err, "couldn't unmarshal request body from json")
		l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err = h.service.Create(ctx, payment)
	if err != nil {
		ctxLogger := log.FromContext(ctx)
		ctxLogger.WithError(err).Error("couldn't create payment")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", payment.ID)
	w.WriteHeader(http.StatusCreated)
	l.Debug("posting completed")
}

func (h *PaymentHandler) UpdateThirdParty(party string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		id := chi.URLParam(r, "id")
		fields := logrus.Fields{}
		fields["method"] = []string{"put"}
		fields["url"] = []string{r.URL.String()}
		fields["party"] = []string{party}
		fields["id"] = []string{id}
		l := h.logger.WithFields(fields)
		ctx := log.ContextWithLogger(r.Context(), l)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = errors.Wrap(err, "couldn't read request body.")
			l.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var thirdParty app.ThirdParty
		if err := json.Unmarshal(body, &thirdParty); err != nil {
			err = errors.Wrap(err, "couldn't unmarshal request body from json")
			l.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		thirdParty.PaymentID = id

		if err := h.service.UpdateThirdParty(ctx, thirdParty, party); err != nil {
			ctxLogger := log.FromContext(ctx)
			ctxLogger.WithError(err).Error("couldn't update third party")
			if errors.Cause(err) == app.ErrNoRows {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		l.Debug("updating completed")
	}
}
