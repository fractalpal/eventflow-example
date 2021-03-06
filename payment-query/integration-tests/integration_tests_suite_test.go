// +build integration query

package integration_tests_test

import (
	"context"
	"sync"
	"testing"

	payment_query "github.com/fractalpal/eventflow-example/payment-query"

	"github.com/fractalpal/eventflow"
	"github.com/fractalpal/eventflow-example/api/http"
	"github.com/fractalpal/eventflow-example/payment"
	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var (
	wg     *sync.WaitGroup
	server *http.Server
	config payment.Config
)

var _ = BeforeEach(func() {
	wg.Add(1)
})

var _ = AfterEach(func() {
	wg.Done()
})

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var _ = BeforeSuite(func() {
	wg = &sync.WaitGroup{}
	Expect(envconfig.Process("payment", &config)).ToNot(HaveOccurred())

	server = payment_query.Initialize(context.Background(), logrus.New(), eventflow.InMemory())
	go server.Start()

})

var _ = AfterSuite(func() {
	go func() {
		wg.Wait()
		_ = server.Shutdown(context.Background())
	}()
})
