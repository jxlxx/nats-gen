package micro

type Testing struct {
	Package      string
	Imports      []string
	MockHandlers []MockHandler
}

type Tests struct {
	TestName    string
	Description string
}

type MockHandler struct {
	OperationID string
	HandlerArgs string
}

type OptionValue struct {
	Name  string
	Value string
}

const TestTemplate = `
package {{ .Package }}

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	{{ range .Imports }}
	{{ . }}
	{{ end }}
)

type MockHandler struct{}

{{ range .MockHandlers }}
func (m MockHandler) {{ .OperationID }}(r micro.Request{{.HandlerArgs}}) {
}
{{ end }}
	
func Test{{.TestName}}(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("err: nats connection, %s\n", err)
	}
	h := MockHandler{}
	if _, err := CreateService(nc, h, Options{{{ range .InitialValues }}
	{{ .Name }}: {{ .Value }}
	{{ end }}
	}); err != nil {
		t.Fatalf("err creating service: %s\n", err)
	}
}
`
