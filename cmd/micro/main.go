package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/jxlxx/nats-gen/micro"
)

var (
	outputFile  string
	inputFile   string
	packageName string
)

func init() {
	flag.StringVarP(&inputFile, "input", "i", "", "help message for flagname")
	flag.StringVarP(&outputFile, "output", "o", "", "help message for flagname")
	flag.StringVarP(&packageName, "package", "p", "", "help message for flagname")

}

func main() {
	flag.Parse()

	var m micro.Spec
	read(inputFile, &m)
	c := micro.CodeGen{
		Package: packageName,
	}

	addedParams := map[string]bool{}
	for _, g := range m.Groups {
		subject, params := parseGroup(g)
		c.Groups = append(c.Groups, micro.ServiceGroup{
			Name:    g.Name,
			Subject: subject,
			Params:  params,
		})
		for _, p := range g.Parameters {
			cc := toCamelCase(p.Name)
			_, ok := addedParams[cc]
			if ok {
				continue
			}
			c.Options = append(c.Options, micro.Option{
				Name: cc,
				Type: p.Type,
			})
			addedParams[cc] = true
		}
	}

	for _, e := range m.Endpoints {
		c.Endpoints = append(c.Endpoints, micro.ServiceEndpoint{
			Group:       e.Group,
			Subject:     e.Subject,
			OperationID: e.OperationID,
		})
		c.WrapperFunctions = append(c.WrapperFunctions, micro.Function{
			OperationID:     e.OperationID,
			MethodSignature: methodSignature(e.OperationID, e.Parameters),
			HandlerArgs:     handlerArgs(e.Parameters),
		})
		if e.Payload.Name != "" {

		}
	}

	if err := writeService(c); err != nil {
		fmt.Println(err)
	}

}

func methodSignature(name string, params []micro.Param) string {
	if len(params) == 0 {
		return fmt.Sprintf("%s(micro.Request)", name)
	}
	types := []string{}
	for _, p := range params {
		types = append(types, p.Type)
	}
	return fmt.Sprintf("%s(micro.Request, %s)", name, strings.Join(types, ", "))
}

func handlerArgs(params []micro.Param) string {
	if len(params) == 0 {
		return "(r)"
	}
	types := []string{}
	for _, p := range params {

		types = append(types, p.Type)
	}
	return fmt.Sprintf("(r, %s)", strings.Join(types, ", "))
}

func parseGroup(g micro.Group) (string, []string) {
	paramNames := []string{}
	s := g.Subject
	for _, p := range g.Parameters {
		cc := toCamelCase(p.Name)
		paramNames = append(paramNames, cc)
		s = strings.ReplaceAll(s, fmt.Sprintf("{%s}", p.Name), "%s")
	}
	return s, paramNames
}

func parseSubject(s string, p map[string]micro.Param) (string, []string) {
	return "", nil
}

func writeService(spec micro.CodeGen) error {
	tmpl := template.Must(template.New("micro-gen").Parse(micro.Template))
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(file, spec); err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}
	return nil
}

func read(filename string, conf interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("err getting working dir: ", err)
	}
	fullPath := fmt.Sprintf("%s/%s", wd, filename)
	f, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalln("err reading yaml: ", err)
	}
	err = yaml.Unmarshal(f, conf)
	if err != nil {
		log.Fatalln("err unmarshal: ", err)
	}
}

func toCamelCase(s string) string {
	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	prevIsCap := false
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		} else if prevIsCap && vIsCap {
			v += 'a'
			v -= 'A'
		}
		prevIsCap = vIsCap

		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
