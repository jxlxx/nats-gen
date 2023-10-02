package bank

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type MockHandler struct{}

func (m MockHandler) NewAccount(r micro.Request, ownerID uuid.UUID) {}
func (m MockHandler) Account(r micro.Request, accountID uuid.UUID)  {}
func (m MockHandler) Accounts(r micro.Request, ownerID uuid.UUID)   {}
func (m MockHandler) Deposit(r micro.Request, deposit Deposit)      {}
func (m MockHandler) Transfer(r micro.Request, transfer Transfer)   {}
func (m MockHandler) Hold(r micro.Request, hold Hold)               {}

func TestCreateService(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("err: nats connection, %s\n", err)
	}
	h := MockHandler{}
	if _, err := CreateService(nc, h, &Options{
		BankCode:    "BMO",
		CountryCode: "CAN",
	}); err != nil {
		t.Fatalf("err creating service: %s\n", err)
	}
}
