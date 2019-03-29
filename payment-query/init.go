package payment_query

import (
	"context"
	"time"

	"github.com/fractalpal/eventflow"
	"github.com/kelseyhightower/envconfig"

	"github.com/fractalpal/eventflow-example/api/http"
	"github.com/fractalpal/eventflow-example/payment-query/adapters/aggregator"
	"github.com/fractalpal/eventflow-example/payment-query/adapters/repository"
	"github.com/fractalpal/eventflow-example/payment-query/adapters/service"
	"github.com/fractalpal/eventflow-example/payment-query/adapters/store"
	queryHttp "github.com/fractalpal/eventflow-example/payment-query/api/http"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context, l logrus.FieldLogger, subscriber eventflow.Subscriber) *http.Server {
	// read env vars
	type Config struct {
		ListenHost string `envconfig:"query_listen_host" default:"localhost"`
		ListenPort string `envconfig:"query_listen_port" default:"8081"`
		// we could split user and password from below
		MongoURL string `envconfig:"query_mongo_url" default:"mongodb://usr:pwd@localhost:27017"`
	}

	var config Config
	err := envconfig.Process("query", &config)

	collection, err := store.NewMongoCollection(ctx,
		config.MongoURL,
		"query",
		"payments",
		time.Second*10,
	)
	if err != nil {
		l.WithError(err).Fatal("couldn't connect to mongo")
	}
	//collection.Indexes() we should have index on 'id'
	queryRepo := repository.NewBasic(collection)
	queryService := service.NewQuery(queryRepo)
	// http server
	timeout := time.Minute
	addr := config.ListenHost + ":" + config.ListenPort
	fields := logrus.Fields{}
	fields["server"] = []string{"payment query"}
	fields["addr"] = []string{addr}
	l = l.WithFields(fields)
	queryRouter := defaultRouter(l.WithField("router", "query"))
	queryServer := http.NewServer(
		http.ServerAddr(addr),
		http.ServerLogger(l),
		http.ServerRouter(queryRouter),
		http.ServerReadTimout(timeout),
		http.ServerWriteTimeout(timeout),
	)

	fields["api"] = []string{"v1"}
	fields["transport"] = []string{"http"}
	queryHandler := queryHttp.NewHandler(l.WithFields(fields), queryService)
	queryServer.Get("/{id}", queryHandler.Get)
	queryServer.Get("/", queryHandler.List)

	// aggregator
	aggr := aggregator.NewMemory(ctx, queryRepo)
	subscriber.Subscribe(aggr, app.PaymentCreated, app.PaymentDeleted, app.BeneficiaryUpdated, app.DebtorUpdated)

	return queryServer
}

func defaultRouter(l logrus.FieldLogger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(http.MaxBytesReader(1024))
	router.Use(middleware.NoCache)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("content-type", "application/json"))
	router.Use(middleware.Timeout(time.Second * 60))
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: l, NoColor: false}))
	return router
}
