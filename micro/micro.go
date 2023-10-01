package micro

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jxlxx/nats-gen/micro/spec"
)

type Generator struct {
	microservices []*Microservice
}

const (
	PACKAGE_UUID string = "github.com/google/uuid"
)

func New() *Generator {
	return &Generator{}
}

//go:embed templates/*
var templatesFS embed.FS

func (g *Generator) Write() error {
	tmpl := template.New("micro")
	t, err := tmpl.ParseFS(templatesFS, "templates/*")
	if err != nil {
		return err
	}
	for _, m := range g.microservices {
		file, err := os.Create(m.File)
		if err != nil {
			return err
		}
		if err := t.ExecuteTemplate(file, "Main", m); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Run(s spec.Spec) ([]string, error) {
	microservices, err := g.GenerateMicroservices(s)
	if err != nil {
		return nil, err
	}
	kvs := []string{}
	return append(microservices, kvs...), nil
}

func (g *Generator) GenerateMicroservices(s spec.Spec) ([]string, error) {
	output := []string{}
	for _, spec := range s.Microservices {
		m := &Microservice{
			Imports:        map[string]Import{},
			File:           spec.TargetFile,
			Package:        spec.Package,
			InitParameters: map[string]Parameter{},
			groupMap:       map[string]Group{},
		}
		if err := m.ParseTypes(spec.Schemas); err != nil {
			return nil, err
		}
		if err := m.ParseGroups(spec.Groups); err != nil {
			return nil, err
		}
		if err := m.ParseEndpoints(spec.Endpoints); err != nil {
			return nil, err
		}
		g.microservices = append(g.microservices, m)
		output = append(output, m.File)
	}
	return output, nil
}

func (m *Microservice) ParseTypes(schemas []spec.Schema) error {
	typeNames := map[string]bool{}
	for _, schema := range schemas {
		if _, ok := typeNames[schema.Name]; ok {
			return fmt.Errorf("duplicate type names, two instances of %s", schema.Name)
		}

		fields := []Field{}
		fieldNames := map[string]bool{}
		for _, v := range schema.Fields {
			if _, ok := fieldNames[v.Name]; ok {
				return fmt.Errorf("duplicate field name in %s schema, two instances of %s", schema.Name, v.Name)
			}
			dataType, err := m.getDataType(v)
			if err != nil {
				return err
			}
			f := Field{
				Name:        v.Name,
				Description: v.Description,
				DataType:    dataType,
			}
			fieldNames[v.Name] = true
			fields = append(fields, f)
		}

		typeNames[schema.Name] = true

		m.Types = append(m.Types, NewType{
			Name:        schema.Name,
			Description: schema.Description,
			Fields:      fields,
		})
	}
	return nil
}

func (m *Microservice) getDataType(field spec.Value) (string, error) {

	// is it a newly defined type ? (schema)
	if field.Schema != "" {
		// TODO check if type exists ? hmmm
		return field.Schema, nil
	}

	// is it a type with special formatting (uuid)
	if field.Format == "uuid" && field.Type == "string" && field.Schema == "" {
		m.Imports["uuid"] = Import{
			URL:  PACKAGE_UUID,
			Name: "uuid",
		}
		return "uuid.UUID", nil
	}

	// is it a base type? (type)
	if isBaseType(field.Type) {

	}

	return field.Type, nil
}

var baseTypes = map[string]bool{
	"string":  true,
	"int":     true,
	"array":   true,
	"float64": true,
}

func isBaseType(t string) bool {
	return baseTypes[t]
}

func (m *Microservice) ParseGroups(groups []spec.Group) error {
	for _, g := range groups {
		subject, args := parseGroupSubject(g)
		if g.Name == "" {
			return fmt.Errorf("cannot have empty group name")
		}
		for _, p := range args {
			m.InitParameters[p.Name] = Parameter{Name: p.Name, DataType: p.DataType}
		}
		group := Group{
			Name:        g.Name,
			Description: g.Description,
			Subject:     subject,
			SubjectArgs: args,
		}
		m.groupMap[g.Name] = group
		m.Groups = append(m.Groups, group)
	}
	return nil
}

func parseGroupSubject(g spec.Group) (string, []Argument) {
	if len(g.Subject.Arguments) == 0 {
		return strings.Join(g.Subject.Tokens, "."), nil
	}
	args := []Argument{}
	s := ""

	argMap := map[string]spec.Value{}
	for _, a := range g.Subject.Arguments {
		argMap[a.Name] = a
	}
	for _, t := range g.Subject.Tokens {
		s = fmt.Sprintf("%s.%s", s, t)
		arg, ok := argMap[t]
		if !ok {
			continue
		}
		args = append(args, Argument{
			Name:     arg.Name,
			Value:    fmt.Sprintf("%s%s", g.Name, arg.Name),
			DataType: arg.Type, // TODO: uuid format
		})
	}
	return s, args

}

func (m *Microservice) ParseEndpoints(endpoints []spec.Endpoint) error {
	for _, e := range endpoints {
		groupName := "service"
		if e.Group != "" {
			g, ok := m.groupMap[e.Group]
			if !ok {
				return fmt.Errorf("undefined group: %s", e.Group)
			}
			groupName = g.Name
		}
		params := []Parameter{}
		args := []Argument{}
		for _, p := range e.Subject.Parameters {
			dataType, err := m.getDataType(p)
			if err != nil {
				return err
			}
			params = append(params, Parameter{
				Name:     p.Name,
				DataType: dataType,
			})
			args = append(args, Argument{
				Name:     p.Name,
				DataType: dataType,
				Value:    p.Name,
			})
		}
		handler := Handler{
			Parameters: params,
			Arguments:  args,
		}
		subject := Subject{}
		m.Endpoints = append(m.Endpoints, Endpoint{
			Name:        e.Name,
			OperationID: e.OperationID,
			Handler:     handler,
			Subject:     subject,
			Group:       groupName,
		})
	}
	return nil
}

func (g *Generator) ParseKVs() error {
	return nil
}
