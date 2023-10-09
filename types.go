package natsgen

type Microservice struct {
	Package        string
	Config         Config
	File           string
	Testing        Testing
	Imports        map[string]Import
	InitParameters map[string]Parameter
	Types          []NewType
	Enums          []Enum
	Groups         []Group
	Endpoints      []Endpoint

	groupMap map[string]Group
	typeMap  map[string]NewType
	enumMap  map[string]map[string]bool
}

type Testing struct {
	File    string
	Name    string
	Package string
	Options map[string]string
}

type Config struct {
	Name        string
	Version     string
	Description string
	QueueGroup  string
	Metadata    map[string]string
}

type Group struct {
	Name        string
	Description string
	Subject     string
	SubjectArgs []Argument
}

type Enum struct {
	Description string
	Name        string
	Values      []EnumValue
}

type EnumValue struct {
	Name  string
	Type  string
	Value string
}

type Endpoint struct {
	Name        string
	QueueGroup  string
	Description string
	OperationID string
	Group       string
	Subject
	Payload
	Handler
}

type Import struct {
	URL        string
	Name       string
	CustomName string
}

type NewType struct {
	Name        string
	Description string
	Fields      []Field
}

type Handler struct {
	Name         string
	Description  string
	OperationID  string
	Parameters   []Parameter
	Arguments    []Argument
	ReturnTypes  []string
	ReturnValues []string

	InterfaceParams string
	HandlerCallArgs string
}

type Payload struct {
	Name        string
	Deserialize bool
	Type        string
	Fields      []Field
}
type Subject struct {
	Name        string
	Description string
	Deserialize bool

	Tokens        []string
	TokenIndexMap map[string]int

	Template   string
	Parameters []Parameter
}

type Parameter struct {
	Name                 string
	Description          string
	DataType             string
	TokenIndex           int
	OnErrorHandlerReturn string
}

type Field struct {
	Name        string
	Description string
	DataType    string
	Tags        string
}

type Argument struct {
	Name        string
	Description string
	Value       string
	DataType    string
}
