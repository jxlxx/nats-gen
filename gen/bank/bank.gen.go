package bank

// This a preamble
// I'm gonna say something like: DO NOT MODIFY!!!

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Service interface {
	NewAccount(r micro.Request, ownerID uuid.UUID)
	Account(r micro.Request, ownerID uuid.UUID, accountID uuid.UUID)
	Accounts(r micro.Request, ownerID uuid.UUID)
	Deposit(r micro.Request, deposit Deposit)
	Transfer(r micro.Request, transfer Transfer)
	Hold(r micro.Request, hold Hold)
}

type Options struct {
	micro.Config
	Name    string
	Version string

	BankCode    string
	CountryCode string
}

// TestEnum
// This is a desciprtion about this enum
type TestEnum string

const (
	Hello   TestEnum = "Hello"
	Gootbye TestEnum = "Gootbye"
)

type Funds struct {
	CurrencyCode string
	Total        int
}

type Account struct {
	ID    uuid.UUID
	funds Funds
}

type Deposit struct {
	ID    uuid.UUID
	Funds int
}

type Transfer struct {
	ID    uuid.UUID
	Funds int
}

type Hold struct {
	ID    uuid.UUID
	Funds int
}

type ServiceWrapper struct {
	Handler Service
}

func CreateService(nc *nats.Conn, s Service, opts *Options) (micro.Service, error) {
	conf, err := createConfig(opts)
	if err != nil {
		return nil, err
	}
	service, err := micro.AddService(nc, conf)
	if err != nil {
		return nil, err
	}
	sw := ServiceWrapper{
		Handler: s,
	}

	base := service.AddGroup(fmt.Sprintf("bank.%s.%s", opts.CountryCode, opts.BankCode))
	admin := service.AddGroup(fmt.Sprintf("admin.bank.%s.%s", opts.CountryCode, opts.BankCode))

	if err := base.AddEndpoint("new", micro.HandlerFunc(sw.NewAccount),
		micro.WithEndpointSubject("new.*")); err != nil {
		return nil, err
	}
	if err := base.AddEndpoint("account", micro.HandlerFunc(sw.Account),
		micro.WithEndpointSubject("account.*.*")); err != nil {
		return nil, err
	}
	if err := base.AddEndpoint("accounts", micro.HandlerFunc(sw.Accounts),
		micro.WithEndpointSubject("accounts.*")); err != nil {
		return nil, err
	}
	if err := admin.AddEndpoint("deposit", micro.HandlerFunc(sw.Deposit),
		micro.WithEndpointSubject("deposit")); err != nil {
		return nil, err
	}
	if err := admin.AddEndpoint("transfer", micro.HandlerFunc(sw.Transfer),
		micro.WithEndpointSubject("transfer")); err != nil {
		return nil, err
	}
	if err := admin.AddEndpoint("hold", micro.HandlerFunc(sw.Hold),
		micro.WithEndpointSubject("hold")); err != nil {
		return nil, err
	}

	return service, nil
}

func createConfig(opts *Options) (micro.Config, error) {
	return micro.Config{
		Name:        opts.Name,
		Version:     opts.Version,
		Description: opts.Description,
	}, nil
}

func (s *ServiceWrapper) NewAccount(r micro.Request) {
	ownerID, err := deserializeNewAccountSubject(r.Subject())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.NewAccount(r, ownerID)
}

func deserializeNewAccountSubject(subj string) (uuid.UUID, error) {
	tokens := strings.Split(subj, ".")

	ownerID, err := uuid.Parse(tokens[1])
	if err != nil {
		return uuid.Nil, err
	}

	return ownerID, nil
}

func (s *ServiceWrapper) Account(r micro.Request) {
	ownerID, accountID, err := deserializeAccountSubject(r.Subject())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.Account(r, ownerID, accountID)
}

func deserializeAccountSubject(subj string) (uuid.UUID, uuid.UUID, error) {
	tokens := strings.Split(subj, ".")

	ownerID, err := uuid.Parse(tokens[1])
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	accountID, err := uuid.Parse(tokens[2])
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return ownerID, accountID, nil
}

func (s *ServiceWrapper) Accounts(r micro.Request) {
	ownerID, err := deserializeAccountsSubject(r.Subject())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.Accounts(r, ownerID)
}

func deserializeAccountsSubject(subj string) (uuid.UUID, error) {
	tokens := strings.Split(subj, ".")

	ownerID, err := uuid.Parse(tokens[1])
	if err != nil {
		return uuid.Nil, err
	}

	return ownerID, nil
}

func (s *ServiceWrapper) Deposit(r micro.Request) {
	deposit, err := deserializeDepositPayload(r.Data())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.Deposit(r, deposit)
}

func deserializeDepositPayload(b []byte) (Deposit, error) {
	p := Deposit{}
	err := json.Unmarshal(b, &p)
	return p, err
}

func (s *ServiceWrapper) Transfer(r micro.Request) {
	transfer, err := deserializeTransferPayload(r.Data())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.Transfer(r, transfer)
}

func deserializeTransferPayload(b []byte) (Transfer, error) {
	p := Transfer{}
	err := json.Unmarshal(b, &p)
	return p, err
}

func (s *ServiceWrapper) Hold(r micro.Request) {
	hold, err := deserializeHoldPayload(r.Data())
	if err != nil {
		if err := r.Error("code", "description", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.Hold(r, hold)
}

func deserializeHoldPayload(b []byte) (Hold, error) {
	p := Hold{}
	err := json.Unmarshal(b, &p)
	return p, err
}
