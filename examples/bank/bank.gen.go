		
package bank

import (
	"fmt"
	
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	
)

type Handler interface {  
	NewAccount(micro.Request, string) 
	Account(micro.Request, string) 
	Accounts(micro.Request, string) 
	Deposit(micro.Request) 
	Transfer(micro.Request) 
	Hold(micro.Request)
}

type ServiceWrapper struct {
	Handler Handler
}

type Options struct {
	Name        string
	Version     string
	Description string 
	CountryCode string 
	BankCode string 
}

func CreateService(nc *nats.Conn, h Handler, opts Options) (micro.Service, error) {
	conf := micro.Config{
		Name:        opts.Name,
		Version:     opts.Version,
		Description: opts.Description,
	}
	service, err := micro.AddService(nc, conf)
	if err != nil {
		return nil, err
	}
	s := ServiceWrapper{
		Handler: h,
	}
	
	base := service.AddGroup(fmt.Sprintf("bank.%s.%s", opts.CountryCode , opts.BankCode ))
	admin := service.AddGroup(fmt.Sprintf("admin.bank.%s.%s", opts.CountryCode , opts.BankCode ))
	
	if err := base.AddEndpoint("new.{id}", micro.HandlerFunc(s.NewAccount)); err != nil {		
		return nil, err
	}
	if err := base.AddEndpoint("account.{id}", micro.HandlerFunc(s.Account)); err != nil {		
		return nil, err
	}
	if err := base.AddEndpoint("accounts.{id}", micro.HandlerFunc(s.Accounts)); err != nil {		
		return nil, err
	}
	if err := admin.AddEndpoint("deposit", micro.HandlerFunc(s.Deposit)); err != nil {		
		return nil, err
	}
	if err := admin.AddEndpoint("transfer", micro.HandlerFunc(s.Transfer)); err != nil {		
		return nil, err
	}
	if err := admin.AddEndpoint("hold", micro.HandlerFunc(s.Hold)); err != nil {		
		return nil, err
	}
	
	return service, nil
}


func (s *ServiceWrapper) NewAccount(r micro.Request) {
	s.Handler.NewAccount(r, string)
}

func (s *ServiceWrapper) Account(r micro.Request) {
	s.Handler.Account(r, string)
}

func (s *ServiceWrapper) Accounts(r micro.Request) {
	s.Handler.Accounts(r, string)
}

func (s *ServiceWrapper) Deposit(r micro.Request) {
	s.Handler.Deposit(r)
}

func (s *ServiceWrapper) Transfer(r micro.Request) {
	s.Handler.Transfer(r)
}

func (s *ServiceWrapper) Hold(r micro.Request) {
	s.Handler.Hold(r)
}


	