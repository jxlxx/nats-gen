package bank

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/nats"
)

type MockHandler struct{}

func (m MockHandler) NewAccount(r micro.Request, ownerID uuid.UUID)                   {}
func (m MockHandler) Account(r micro.Request, ownerID uuid.UUID, accountID uuid.UUID) {}
func (m MockHandler) Accounts(r micro.Request, ownerID uuid.UUID)                     {}
func (m MockHandler) Deposit(r micro.Request, deposit Deposit)                        {}
func (m MockHandler) Transfer(r micro.Request, transfer Transfer)                     {}
func (m MockHandler) Hold(r micro.Request, hold Hold)                                 {}

func TestBankingService(t *testing.T) {
	ctx := context.Background()

	natsContainer, err := tc.RunContainer(ctx,
		testcontainers.WithImage("nats:2"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Clean up the container
	defer func() {
		if err := natsContainer.Terminate(ctx); err != nil {
			t.Fatalf(err.Error())
		}
	}()

	connectionURL, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	nc, err := nats.Connect(connectionURL)
	if err != nil {
		t.Fatalf(err.Error())
	}

	h := MockHandler{}
	if _, err := CreateService(nc, h, &Options{
		BankCode:    "BMO",
		CountryCode: "CAN",
		Name:        "BankingService",
		Version:     "0.0.1",
	}); err != nil {
		t.Fatalf("err creating service: %s\n", err)
	}
}
