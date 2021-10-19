package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strings"
)

var packageName string

var output = &bytes.Buffer{}

func main() {
	fset := token.NewFileSet()
	file := os.Args[1]
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	packageName = f.Name.Name
	fmt.Fprintln(output, "package "+packageName)
	fmt.Fprint(output, `import (
	"go.bug.st/json"
	"errors"
)
`)
	for _, c := range f.Comments {
		comment := c.Text()
		if strings.HasPrefix(comment, "lsp:generate ") {
			generateType(strings.TrimSpace(comment[13:]))
		}
	}

	formatted, err := format.Source(output.Bytes())
	if err != nil {
		fmt.Println("Error formatting output code:", err)
		os.Exit(1)
	}
	os.Stdout.Write(formatted)
}

func generateType(cmdline string) {
	cmd := strings.Split(cmdline, " ")
	if len(cmd) != 3 {
		fmt.Println(`syntax error: lsp:generate should be in the format "TYPE!TYPE!TYPE as NEWTYPE"`)
		os.Exit(1)
	}
	if cmd[1] != "as" {
		fmt.Println("syntax error: expected 'as' keyword")
		os.Exit(1)
	}
	type newtype struct {
		Name    string
		Types   []string
		HasNull bool
	}
	t := newtype{
		Name:  cmd[2],
		Types: strings.Split(cmd[0], "|"),
	}
	for i, n := range t.Types {
		if n == "Null" {
			t.Types = append(t.Types[:i], t.Types[i+1:]...)
			t.HasNull = true
			break
		}
	}
	funcs := template.FuncMap{"join": strings.Join}
	templ := template.Must(template.New("class").Funcs(funcs).Parse(`
// {{.Name}} is an intersection type of {{join .Types ", "}}
type {{.Name}} []byte

// MarshalJSON implements json.Marshaler interface
func (m {{.Name}}) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *{{.Name}}) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	if m.Value() == nil {
		// return errors.New("invalid input: expected one of: string, int, bool, Array, Object, Null")
		return errors.New("invalid input: expected one of: {{join .Types ", "}}")
	}
	return nil
}

// Value returns the value of the sum-type
func (r *{{.Name}}) Value() interface{} {
	data := []byte(*r)
{{if .HasNull}}
	var n Null
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}
{{end}}
{{range $i, $t := .Types}}
	var r{{$i}} {{$t}}
	if err := json.Unmarshal(data, &r{{$i}}); err == nil {
		return r{{$i}}
	}
{{end}}

	return nil
}

`))

	if err := templ.Execute(output, t); err != nil {
		fmt.Printf("Error executing template: %s", err)
		os.Exit(2)
	}
}
