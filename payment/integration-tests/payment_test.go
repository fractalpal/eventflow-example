package integration_tests_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

func createPaymentRequest(url, json string) (*http.Response, error) {
	return http.Post(url, "application/json", bytes.NewBufferString(json))
}

var _ = Describe("Payment", func() {

	// default client for tests
	client := http.Client{
		Timeout: time.Second * 30,
	}

	json := `
	{
	"type":"instant",
	"attributes":
		{"amount":"100.50",
		"currency":"EUR",
		"beneficiary_party":
			{"account_name":"Ben",
			"account_number":"123"
			},
		"debtor_party":
			{"account_name":"Deb",
			"account_number":"987"},
		"payment_id":"92aaf311-a2fe-4022-86ef-162f314149df",
		"payment_type":"credit_card",
		"processing_date":"2019-03-13T21:20:57+01:00",
		"reference":"f1eef151-87c0-4aac-afbc-3e68c8f5807c"}}
	`

	var baseURL string
	BeforeEach(func() {
		baseURL = fmt.Sprintf("http://%s:%s/", config.ListenHost, config.ListenPort)
	})

	It("Returns location of newly created payment", func() {
		// given
		requestBody := json

		// when
		resp, err := createPaymentRequest(baseURL, requestBody)

		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		Expect(resp.Header.Get("Location")).ToNot(BeEmpty())
	})

	It("Deletes payment", func() {
		// given
		resp, err := createPaymentRequest(baseURL, json)
		Expect(err).ToNot(HaveOccurred())
		id := resp.Header.Get("Location")

		// when
		req, err := http.NewRequest(http.MethodDelete, baseURL+id, http.NoBody)
		Expect(err).ToNot(HaveOccurred())
		resp, err = client.Do(req)

		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

})
