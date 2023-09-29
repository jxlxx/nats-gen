package micro

const Template = `		
package {{ .Package }}

import (
	"fmt"
	
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	{{ range .Imports }}
	{{ . }}
	{{ end }}
)

type Handler interface { {{ range .WrapperFunctions }} 
	{{ .MethodSignature }}{{ end }}
}

type ServiceWrapper struct {
	Handler Handler
}

type Options struct {
	Name        string
	Version     string
	Description string {{ range .Options }}
	{{ .Name  }} {{ .Type }} {{ end }}
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
	{{ range .Groups }}
	{{ .Name }} := service.AddGroup(fmt.Sprintf("{{ .Subject}}"{{ range .Params }}, opts.{{.}} {{end}})){{ end }}
	{{ range .Endpoints }}
	if err := {{ .Group }}.AddEndpoint("{{ .Subject }}", micro.HandlerFunc(s.{{ .OperationID }})); err != nil {		
		return nil, err
	}{{ end }}
	
	return service, nil
}

{{ range .WrapperFunctions }}
func (s *ServiceWrapper) {{ .OperationID }}(r micro.Request) {
	s.Handler.{{ .OperationID }}{{ .HandlerArgs }}
}
{{ end }}

	`
