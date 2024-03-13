package goveikkaus

import (
	"context"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/j-flat/go-veikkaus/internal/veikkausapi"
)

const HappyCaseBalance = 1577

var _ = Describe("authservice: account balance", func() {
	var client *Client
	var mux *http.ServeMux
	var teardown func()

	var happyCaseResponseBytes = []byte(`{"status":"ACTIVE","timerInterval":60,"balances":{"CASH":{"currency":"EUR","type":"CASH","balance":1577,"usableBalance":1577,"frozenBalance":0,"holdBalance":0}}}`)
	var unknownErrorBytes = []byte(`{"code": "UNKNOWN", "fieldErrors": []}`)

	BeforeEach(func() {
		client, mux, _, teardown = setup()
	})

	AfterEach(func() {
		defer teardown()
	})

	DescribeTable("AccountBalance",
		func(shouldSucceed bool, expectedStatusCode, expectedBalance int, expectedErr error, expectedResponseBodyBytes []byte) {
			mux.HandleFunc("/"+api.AccountBalanceEndpoint, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(expectedStatusCode)
				if _, err := w.Write(happyCaseResponseBytes); err != nil {
					log.Fatalf("Error whilte writing the response body in unit-test: %v", err)
				}
			})

			ctx := context.Background()
			data, _, err := client.Auth.AccountBalance(ctx)

			if shouldSucceed {
				Expect(data).To(BeAssignableToTypeOf(&AccountBalance{}))
				Expect(data.Balances.Cash.Balance).To(Equal(expectedBalance))
				Expect(err).To(BeNil())
			} else {
				Expect(err).To(BeAssignableToTypeOf(expectedErr))
				Expect(data).To(BeNil())
			}
		},
		Entry("should return account balance on happy-case", true, http.StatusOK, HappyCaseBalance, nil, happyCaseResponseBytes),
		Entry("should return error when response status code is unsupported, but response is otherwise successful", false, http.StatusMovedPermanently, nil, &api.UnsupportedStatusCodeError{}, happyCaseResponseBytes),
		Entry("should return error when response errored and code is unknown", false, http.StatusBadRequest, nil, &api.APIErrorNotImplementedError{}, unknownErrorBytes),
	)
})
