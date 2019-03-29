package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fractalpal/eventflow"
	"github.com/fractalpal/eventflow-example/api/http"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment"
	query "github.com/fractalpal/eventflow-example/payment-query"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ListerAddr string
}

func main() {
	// Application Context.
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// logger
	l := log.NewLogrus(logrus.DebugLevel).WithField("system", []string{"cqrs/es"})
	appCtx = log.ContextWithLogger(appCtx, l)
	//l.SetFormatter(&logrus.JSONFormatter{})

	// flow
	flow := eventflow.InMemory()

	// srv chan
	srvErrChan := make(chan error, 1)

	// query part
	queryServer := query.Initialize(appCtx, l, flow)
	go func() {
		if err := queryServer.Start(); err != nil {
			srvErrChan <- errors.Wrap(err, "query")
		}
	}()

	// payment part
	paymentServer, db := payment.Initialize(l, flow)
	defer db.Close()
	go func() {
		if err := paymentServer.Start(); err != nil {
			srvErrChan <- errors.Wrap(err, "payment")
		}
	}()

	// os.Signal
	osSigChan := osSignals()

	// handle proper channels
	// os.Signals for gracefully shutdown
	// server err
	idleConnsClosed := make(chan struct{})
	go func() {
		for {
			select {
			case sig := <-osSigChan:
				l.Infof("os.Signal notification: '%s', closing server", sig)
				if err := paymentServer.Shutdown(appCtx); err != nil {
					l.Error(errors.Wrap(err, "cannot shutdown payment server"))
				}
				if err := queryServer.Shutdown(appCtx); err != nil {
					l.Error(errors.Wrap(err, "cannot shutdown query server"))
				}
				close(idleConnsClosed)
			case err := <-srvErrChan:
				if errors.Cause(err) == http.ErrServerClosed {
					l.Info(errors.Wrap(err, "server"))
					continue
				}
				l.Error(errors.Wrap(err, "server error"))
			}
		}
	}()
	<-idleConnsClosed
}

func osSignals() chan os.Signal {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	return sigint
}
