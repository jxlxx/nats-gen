package micro

type Microservice struct {
	Package        string
	File           string
	Config         Config
	Imports        map[string]Import
	InitParameters map[string]Parameter
	Types          []NewType
	Groups         []Group
	Endpoints      []Endpoint

	groupMap map[string]Group
}

type Config struct {
	Name        string
	Version     string
	Description string
}

type Group struct {
	Name        string
	Description string
	Subject     string
	SubjectArgs []Argument
}

type Endpoint struct {
	Name        string
	OperationID string
	Group       string
	Subject
	Payload
	Handler
}

type Payload struct {
	Name        string
	Deserialize bool
	Type        string
	Fields      []Field
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
	OperationID  string
	Description  string
	Parameters   []Parameter
	Arguments    []Argument
	ReturnTypes  []string
	ReturnValues []string

	InterfaceParams string
	HandlerCallArgs string
}

type Subject struct {
	Name              string
	Deserialize       bool
	Tokens            []string
	TokenIndexMap     map[string]int
	NumberOfTokens    int
	Description       string
	ExpandedWithGroup string
	Template          string
	Parameters        []Parameter
	Argument          map[string]Argument
}

type Parameter struct {
	Name        string
	Description string
	DataType    string
	TokenIndex  int
	NilValue    string
	OnErrorNils string
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
