package payment

import (
	"database/sql"
	"github.com/kelseyhightower/envconfig"
	"time"

	"github.com/fractalpal/eventflow"
	"github.com/fractalpal/eventflow-example/api/http"
	"github.com/fractalpal/eventflow-example/payment/adapters/commander"
	"github.com/fractalpal/eventflow-example/payment/adapters/repository"
	"github.com/fractalpal/eventflow-example/payment/adapters/service"
	"github.com/fractalpal/eventflow-example/payment/adapters/store"
	"github.com/fractalpal/eventflow-example/payment/adapters/validator"
	paymentHttp "github.com/fractalpal/eventflow-example/payment/api/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
)

// read env vars
type Config struct {
	ListenHost string `envconfig:"payment_listen_host" default:"localhost"`
	ListenPort string `envconfig:"payment_listen_port" default:"8080"`
	// we could split user and password from below
	PostgresURL            string `envconfig:"payment_postgres_url" default:"postgres://usr:pwd@localhost:5432/events?sslmode=disable"`
	PostgresMigrationsPath string `envconfig:"payment_postgres_migrations_path" default:"file://payment/adapters/store/migrations"`
}

func Initialize(l logrus.FieldLogger, publisher eventflow.Publisher) (*http.Server, *sql.DB) {
	var config Config
	err := envconfig.Process("payment", &config)

	// http server log
	l = l.WithField("server", []string{"payment command"})

	// db
	srcUrl := config.PostgresURL
	dbLog := l.WithFields(logrus.Fields{
		"db": []string{srcUrl},
	})
	db, err := store.Postgres(srcUrl)
	if err != nil {
		dbLog.WithError(err).Fatal("couldn't connect to db.")
	}
	if err := db.Ping(); err != nil {
		dbLog.WithError(err).Fatal("couldn't ping to db.")
	}

	m, err := store.Migration(db, config.PostgresMigrationsPath)
	if err != nil {
		l.WithField("migration_path", config.PostgresMigrationsPath).WithError(err).Fatal("couldn't create migration")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		l.WithError(err).Fatal("couldn't migrate up")
	}

	// repo
	repo := repository.NewFlow(
		commander.NewJson(),
		store.New(db),
		publisher,
	)
	// service
	svc := service.New(&service.UUIDProvider{}, repo, validator.NewEmpty())

	addr := config.ListenHost + ":" + config.ListenPort
	fields := logrus.Fields{}
	fields["server"] = []string{"payment command"}
	fields["addr"] = []string{addr}
	l = l.WithFields(fields)
	paymentRouter := defaultRouter(l.WithField("router", "payment"))

	// http server
	timeout := time.Minute
	paymentServer := http.NewServer(
		http.ServerAddr(addr),
		http.ServerLogger(l),
		http.ServerRouter(paymentRouter),
		http.ServerReadTimout(timeout),
		http.ServerWriteTimeout(timeout))
	fields["api"] = []string{"v1"}
	fields["transport"] = []string{"http"}
	handler := paymentHttp.NewHandler(svc, l.WithFields(fields))
	paymentServer.Post("/", handler.Post)
	paymentServer.Put("/{id}/beneficiary-party", handler.UpdateThirdParty("beneficiary"))
	paymentServer.Put("/{id}/debtor-party", handler.UpdateThirdParty("debtor"))
	paymentServer.Delete("/{id}", handler.Delete)

	return paymentServer, db
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
