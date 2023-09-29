package main

type MicroSpec struct {
	Config    Config     `yaml:"config"`
	Groups    []Group    `yaml:"groups"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Config struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Group struct {
	Name       string  `yaml:"name"`
	Subject    string  `yaml:"subject"`
	Parameters []Param `yaml:"parameters"`
}

type Param struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Format string `yaml:"format,omitempty"`
}

type Endpoint struct {
	Name        string  `yaml:"name"`
	Payload     Payload `yaml:"payload,omitempty"`
	Subject     string  `yaml:"subject"`
	OperationID string  `yaml:"operationId"`
	Parameters  []Param `yaml:"parameters,omitempty"`
	Group       string  `yaml:"group"`
}

type Payload struct {
	Name   string  `yaml:"name"`
	Values []Param `yaml:"values"`
}

type CodeGenSpec struct {
	Package          string
	Imports          []string
	Options          []Option
	Groups           []ServiceGroup
	Endpoints        []ServiceEndpoint
	WrapperFunctions []Function
}

type Option struct {
	Name string
	Type string
}

type ServiceGroup struct {
	Name    string
	Subject string
	Params  []string
}

type ServiceEndpoint struct {
	Group       string
	Subject     string
	OperationID string
}

type Function struct {
	OperationID     string
	MethodSignature string
	HandlerArgs     string
}
