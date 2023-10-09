package spec

type Spec struct {
	Microservices []Microservice `yaml:"microservices"`
	KeyValues     []KeyValue     `yaml:"keyValues"`
}

type Microservice struct {
	Package    string  `yaml:"package"`
	TargetFile string  `yaml:"targetFile"`
	Testing    Testing `yaml:"testing"`

	Config    Config     `yaml:"config"`
	Groups    []Group    `yaml:"groups"`
	Endpoints []Endpoint `yaml:"endpoints"`
	Schemas   []Schema   `yaml:"schemas"`
	Enums     []Enum     `yaml:"enums"`
}

type Testing struct {
	Name    string            `yaml:"name"`
	File    string            `yaml:"file"`
	Package string            `yaml:"package"`
	Options map[string]string `yaml:"options"`
	Tests   bool              `yaml:"enable"`
}

type Schema struct {
	Name        string  `yaml:"name"`
	Fields      []Value `yaml:"fields"`
	Description string  `yaml:"description"`
}

type Enum struct {
	Name        string   `yaml:"name"`
	Values      []string `yaml:"values"`
	Description string   `yaml:"description"`
}

type Config struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	QueueGroup  string            `yaml:"queueGroup"`
	Metadata    map[string]string `yaml:"metadata"`
}

type Group struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Subject     Subject `yaml:"subject"`
	QueueGroup  string  `yaml:"queueGroup"`
}

type Endpoint struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	OperationID string            `yaml:"operationId"`
	QueueGroup  string            `yaml:"queueGroup"`
	Group       string            `yaml:"group"`
	Subject     Subject           `yaml:"subject"`
	Payload     Payload           `yaml:"payload"`
	Metadata    map[string]string `yaml:"metadata"`
	Generate    GenerateOptions   `yaml:"generate"`
}

type GenerateOptions struct {
	Tests         bool `yaml:"tests"`
	PayloadSchema bool `default:"true" yaml:"payloadSchema"`
	SubjectSchema bool `default:"true" yaml:"subjectSchema"`
	Documentation bool `default:"true" yaml:"documentation"`
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
	Items       string   `yaml:"items"`
	Enum        []string `yaml:"enum"`
	Format      string   `yaml:"format"`
	Examples    []string `yaml:"examples"`
	Description string   `yaml:"description"`
}

type KeyValue struct{}
