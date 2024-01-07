package cats

// This a preamble
// I'm gonna say something like: DO NOT MODIFY!!!

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Service interface {
	NewCat(r micro.Request, cat CatIntake)
	EditCat(r micro.Request, catID string, cat CatIntake)
	GetCat(r micro.Request, catID string)
}

type Options struct {
	micro.Config
	Name    string
	Version string
}

type CatIntake struct {
	Name       string
	BirthYear  int
	BirthMonth int
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

	cats := service.AddGroup(fmt.Sprintf("cats"))

	if err := cats.AddEndpoint("new", micro.HandlerFunc(sw.NewCat),
		micro.WithEndpointSubject("new")); err != nil {
		return nil, err
	}
	if err := cats.AddEndpoint("edit", micro.HandlerFunc(sw.EditCat),
		micro.WithEndpointSubject("edit.*")); err != nil {
		return nil, err
	}
	if err := cats.AddEndpoint("get", micro.HandlerFunc(sw.GetCat),
		micro.WithEndpointSubject("*")); err != nil {
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

func (s *ServiceWrapper) NewCat(r micro.Request) {
	cat, err := deserializeNewCatPayload(r.Data())
	if err != nil {
		if err := r.Error("err", "payload deserialization error", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.NewCat(r, cat)
}

func deserializeNewCatPayload(b []byte) (CatIntake, error) {
	p := CatIntake{}
	err := json.Unmarshal(b, &p)
	return p, err
}

func (s *ServiceWrapper) EditCat(r micro.Request) {
	catID, err := deserializeEditCatSubject(r.Subject())
	if err != nil {
		if err := r.Error("err", "subject deserialization error", nil); err != nil {
			fmt.Println(err)
		}
		return
	}
	cat, err := deserializeEditCatPayload(r.Data())
	if err != nil {
		if err := r.Error("err", "payload deserialization error", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.EditCat(r, catID, cat)
}

func deserializeEditCatSubject(subj string) (string, error) {
	tokens := strings.Split(subj, ".")

	catID := tokens[1]

	return catID, nil
}

func deserializeEditCatPayload(b []byte) (CatIntake, error) {
	p := CatIntake{}
	err := json.Unmarshal(b, &p)
	return p, err
}

func (s *ServiceWrapper) GetCat(r micro.Request) {
	catID, err := deserializeGetCatSubject(r.Subject())
	if err != nil {
		if err := r.Error("err", "subject deserialization error", nil); err != nil {
			fmt.Println(err)
		}
		return
	}

	s.Handler.GetCat(r, catID)
}

func deserializeGetCatSubject(subj string) (string, error) {
	tokens := strings.Split(subj, ".")

	catID := tokens[0]

	return catID, nil
}
