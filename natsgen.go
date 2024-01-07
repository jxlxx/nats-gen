package natsgen

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/jxlxx/nats-gen/spec"
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

//go:embed templates/micro/*
var microserviceTemplates embed.FS

func (g *Generator) Write() error {
	tmpl := template.New("micro")
	t, err := tmpl.ParseFS(microserviceTemplates, "templates/micro/*")
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
		if m.Testing.File == "" {
			continue
		}
		testFile, err := os.Create(m.Testing.File)
		if err != nil {
			return err
		}
		if err := t.ExecuteTemplate(testFile, "Testing", m); err != nil {
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
			InitParameters: map[string]Parameter{},
			groupMap:       map[string]Group{},
			typeMap:        map[string]NewType{},
			enumMap:        map[string]map[string]bool{},
			File:           spec.TargetFile,
			Package:        spec.Package,
			Testing: Testing{
				Name:    spec.Testing.Name,
				Package: spec.Testing.Package,
				File:    spec.Testing.File,
				Options: spec.Testing.Options,
			},
		}
		if err := m.ParseEnums(spec.Enums); err != nil {
			return nil, err
		}
		fmt.Println("enums parsed")
		if err := m.ParseTypes(spec.Schemas); err != nil {
			return nil, err
		}
		fmt.Println("types parsed")
		if err := m.ParseGroups(spec.Groups); err != nil {
			return nil, err
		}
		fmt.Println("groups parsed")
		if err := m.ParseEndpoints(spec.Endpoints); err != nil {
			return nil, err
		}
		fmt.Println("endpoints parsed")
		g.microservices = append(g.microservices, m)
		output = append(output, m.File, m.Testing.File)
	}
	return output, nil
}

func (m *Microservice) ParseEnums(enums []spec.Enum) error {
	for _, e := range enums {
		if isBaseType(e.Name) {
			return fmt.Errorf("can't used name of basetype as enum: %s", e.Name)
		}
		if _, ok := m.enumMap[e.Name]; ok {
			return fmt.Errorf("already defined this enum: %s", e.Name)
		}
		values := []EnumValue{}
		m.enumMap[e.Name] = map[string]bool{}
		for _, v := range e.Values {
			m.enumMap[e.Name][v] = true
			values = append(values, EnumValue{
				Name:  v,
				Type:  e.Name,
				Value: v,
			})
		}
		enum := Enum{
			Name:        e.Name,
			Description: e.Description,
			Values:      values,
		}
		m.Enums = append(m.Enums, enum)
	}
	return nil
}

func (m *Microservice) ParseTypes(schemas []spec.Schema) error {
	typeNames := map[string]bool{}
	for _, schema := range schemas {
		if _, ok := m.enumMap[schema.Name]; ok {
			return fmt.Errorf("enum name used as schema name %s", schema.Name)
		}
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
				fmt.Println(v)
				fmt.Println(m.typeMap)
				return err
			}

			f := Field{
				Name:        v.Name,
				Description: v.Description,
				DataType:    dataType,
			}
			fields = append(fields, f)
			fieldNames[v.Name] = true
		}

		typeNames[schema.Name] = true
		newType := NewType{
			Name:        schema.Name,
			Description: schema.Description,
			Fields:      fields,
		}
		m.Types = append(m.Types, newType)
		m.typeMap[schema.Name] = newType
	}
	return nil
}

func (m *Microservice) getDataType(field spec.Value) (string, error) {
	if field.Type == "string" && field.Format == "uuid" {
		m.Imports["uuid"] = Import{
			Name: "uuid",
			URL:  PACKAGE_UUID,
		}
		return "uuid.UUID", nil
	}

	if field.Type == "string" && field.Format == "datetime" {
		m.Imports["time"] = Import{
			Name: "time",
			URL:  "time",
		}
		return "time.Time", nil
	}

	if field.Type == "array" && isBaseType(field.Items) {
		return fmt.Sprintf("[]%s", field.Items), nil
	}

	if field.Type == "enum" {
		return field.Schema, nil
	}

	if isBaseType(field.Type) {
		return field.Type, nil
	}

	if _, ok := m.typeMap[field.Schema]; ok {
		return field.Schema, nil
	}

	if _, ok := m.typeMap[field.Type]; ok {
		return field.Type, nil
	}

	return field.Type, fmt.Errorf("unknown type or schema: %s, %s", field.Type, field.Schema)
}

var baseTypes = map[string]bool{
	"string":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"float64": true,
	"float32": true,
	"bool":    true,
	"byte":    true,
}

func isBaseType(t string) bool {
	return baseTypes[t]
}

func (m *Microservice) ParseGroups(groups []spec.Group) error {
	groupNames := map[string]bool{}
	for _, g := range groups {
		if g.Name == "" {
			return fmt.Errorf("cannot have empty group name")
		}
		if _, ok := groupNames[g.Name]; ok {
			return fmt.Errorf("duplicate group name")
		}
		subject, args, err := m.parseGroupSubjectArguments(g)
		if err != nil {
			return err
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

func (m Microservice) parseGroupSubjectArguments(g spec.Group) (string, []Argument, error) {
	if len(g.Subject.Arguments) == 0 {
		return strings.Join(g.Subject.Tokens, "."), nil, nil
	}
	args := []Argument{}
	s := ""

	argMap := map[string]spec.Value{}
	for _, a := range g.Subject.Arguments {
		argMap[a.Name] = a
	}
	for _, t := range g.Subject.Tokens {
		arg, ok := argMap[t]
		if !ok {
			if s == "" {
				s = t
				continue
			}
			s = strings.Join([]string{s, t}, ".")
			continue
		}
		dataType, err := m.getDataType(arg)
		if err != nil {
			return "", nil, err
		}
		s = strings.Join([]string{s, "%s"}, ".")
		args = append(args, Argument{
			Name:     arg.Name,
			Value:    fmt.Sprintf("%s%s", g.Name, arg.Name),
			DataType: dataType, // not all data types are valid here
		})
	}
	return s, args, nil

}

func createTokenIndexMap(tokens []string) map[string]int {
	m := map[string]int{}
	for i, t := range tokens {
		m[t] = i
	}
	return m
}

func replaceParamsWith[T any](tokens []string, params map[string]T, replace string) string {
	result := []string{}
	for _, t := range tokens {
		if _, ok := params[t]; ok {
			t = replace
		}
		result = append(result, t)
	}
	return strings.Join(result, ".")
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
		paramMap := map[string]spec.Value{}
		for _, p := range e.Subject.Parameters {
			paramMap[p.Name] = p
		}
		tokenIndexMap := createTokenIndexMap(e.Subject.Tokens)
		onErrorReturn := m.getErrorReturnString(e.Subject.Tokens, paramMap)

		for _, p := range e.Subject.Parameters {
			dataType, err := m.getDataType(p)
			if err != nil {
				return err
			}
			params = append(params, Parameter{
				Name:                 p.Name,
				Description:          p.Description,
				DataType:             dataType,
				TokenIndex:           tokenIndexMap[p.Name],
				OnErrorHandlerReturn: onErrorReturn,
			})
		}

		tmpl := replaceParamsWith(e.Subject.Tokens, paramMap, "*")
		subject := Subject{
			Template:    tmpl,
			Parameters:  params,
			Deserialize: len(params) > 0,
		}
		payload := Payload{}
		if e.Payload.Schema != "" {
			payload = Payload{
				Name:        e.Payload.Name,
				Type:        e.Payload.Schema,
				Deserialize: true,
			}
		}
		m.Endpoints = append(m.Endpoints, Endpoint{
			Name:        e.Name,
			Description: e.Description,
			OperationID: e.OperationID,
			Group:       groupName,
			Subject:     subject,
			Payload:     payload,
		})
	}
	return nil
}

func (m Microservice) getErrorReturnString(tokens []string, params map[string]spec.Value) string {
	s := []string{}
	for _, t := range tokens {
		if p, ok := params[t]; ok {
			n, err := m.getNilValue(p)
			if err != nil {
				log.Fatal(err)
			}
			s = append(s, n)
		}
	}
	return strings.Join(s, ", ")
}

func (m Microservice) getNilValue(v spec.Value) (string, error) {
	d, err := m.getDataType(v)
	if err != nil {
		return "", err
	}
	if d == "uuid.UUID" {
		return "uuid.Nil", nil
	}
	if d == "time.Time" {
		return "time.Now()", nil
	}
	if v.Type == "array" {
		return "nil", nil
	}

	if v.Type == "enum" {
		return "\"\"", nil
	}
	if isBaseType(d) {
		return baseTypeNilString(d)
	}
	return fmt.Sprintf("%s{}", d), nil
}

func baseTypeNilString(t string) (string, error) {
	switch t {
	case "array":
		return "nil", nil
	case "string":
		return "\"\"", nil
	case "int":
		return "0", nil
	case "int8":
		return "0", nil
	case "int16":
		return "0", nil
	case "int32":
		return "0", nil
	case "int64":
		return "0", nil
	case "uint":
		return "0", nil
	case "uint8":
		return "0", nil
	case "uint16":
		return "0", nil
	case "uint32":
		return "0", nil
	case "uint64":
		return "0", nil
	case "float64":
		return "0", nil
	case "float32":
		return "0", nil
	case "bool":
		return "0", nil
	case "byte":
		return "0", nil
	}
	return "", fmt.Errorf("unknown base type: %s", t)
}

func (g *Generator) GenerateKVs() error {
	return nil
}
