package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// QueryHandler handler
type QueryHandler struct {
	logger       logrus.FieldLogger
	service      app.Service
	defaultPage  int64
	defaultLimit int64
}

// NewHandler returns new QueryHandler
func NewHandler(logger logrus.FieldLogger, service app.Service) *QueryHandler {
	return &QueryHandler{
		logger:       logger,
		service:      service,
		defaultPage:  1,
		defaultLimit: 10, // hard coded for now
	}
}

// Links links info for response
type Links struct {
	Self string `json:"self"`
}

// PaymentResponse http payment response
type PaymentResponse struct {
	Data  app.Payment
	Links Links `json:"links"`
}

// Get handler
func (h *QueryHandler) Get(w http.ResponseWriter, r *http.Request) {
	l := h.logger.WithField("method", []string{"get"})
	ctx := log.ContextWithLogger(r.Context(), l)

	id := chi.URLParam(r, "id")
	payment, err := h.service.FindByID(ctx, id)
	if err != nil {
		if err != app.ErrNoResults {
			l.WithError(err).Error("couldn't found")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		l.Debug("couldn't found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := PaymentResponse{
		Data: payment,
	}
	resp.Links.Self = r.Host + r.URL.Path // build proper self links
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
	l.Debug("payment served")
}

// ListResponse for collections of payments
type ListResponse struct {
	Data  []app.Payment `json:"data"`
	Links Links         `json:"links"`
}

// List handler
func (h *QueryHandler) List(w http.ResponseWriter, r *http.Request) {
	l := h.logger.WithField("method", []string{"get"})
	ctx := log.ContextWithLogger(r.Context(), l)

	query := r.URL.Query()
	pageParam := query.Get("page")
	page, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil || page <= 0 {
		l.WithField("page", pageParam).WithError(err).Debug("incorrect page param, settings default")
		page = h.defaultPage
	}
	limitParam := query.Get("limit")
	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		l.WithField("limit", limitParam).WithError(err).Debug("incorrect limit param, settings default")
		limit = h.defaultLimit
	}

	// filtering and paging not implemented
	list, err := h.service.FindAll(ctx, page, limit)
	if err != nil {
		if err == app.ErrNoResults {
			// we will serve empty list
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var response ListResponse
	response.Data = list
	response.Links.Self = r.Host + r.URL.Path // build proper self links

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
	l.Debug("list served")
}
