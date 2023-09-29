package bank

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type MockHandler struct{}

func (m MockHandler) NewAccount(micro.Request, string) {
}

func (m MockHandler) Account(micro.Request, string) {
}

func (m MockHandler) Accounts(micro.Request, string) {
}

func (m MockHandler) Deposit(micro.Request) {
}

func (m MockHandler) Transfer(micro.Request) {
}

func (m MockHandler) Hold(micro.Request) {
}

func TestBankService(t *testing.T) {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("err: nats connection, %s\n", err)
	}
	h := MockHandler{}
	if _, err := CreateService(nc, h, Options{}); err != nil {
		t.Fatalf("err creating service: %s\n", err)
	}

}
