package spec

type Spec struct {
	Microservices []Microservice `yaml:"microservices"`
	KeyValues     []KeyValue     `yaml:"keyValues"`
}

type Microservice struct {
	Package    string `yaml:"package"`
	TargetFile string `yaml:"targetFile"`
	Tests      bool   `yaml:"tests"`

	Config    Config     `yaml:"config"`
	Groups    []Group    `yaml:"groups"`
	Endpoints []Endpoint `yaml:"endpoints"`
	Schemas   []Schema   `yaml:"schemas"`
}

type Schema struct {
	Name        string  `yaml:"name"`
	Fields      []Value `yaml:"fields"`
	Description string  `yaml:"description"`
}

type Config struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Group struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Subject     Subject `yaml:"subject"`
}

type Endpoint struct {
	Name          string  `yaml:"name"`
	Description   string  `yaml:"description"`
	OperationID   string  `yaml:"operationId"`
	Group         string  `yaml:"group"`
	Subject       Subject `yaml:"subject"`
	PayloadSchema Payload `yaml:"payload"`
}

type Subject struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tokens      []string `yaml:"tokens"`
	Parameters  []Value  `yaml:"parameters"`
	Arguments   []Value  `yaml:"arguments"`
}

type Payload struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Schema      string `yaml:"schema"`
	Type        string `yaml:"type"`
	Format      string `yaml:"format"`
}

type Value struct {
	Name        string   `yaml:"name"`
	Required    bool     `yaml:"required"`
	Type        string   `yaml:"type"`
	Schema      string   `yaml:"schema"`
	Format      string   `yaml:"format"`
	Examples    []string `yaml:"examples"`
	Description string   `yaml:"description"`
}

type KeyValue struct{}
