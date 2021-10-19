package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strings"
)

var packageName string

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	file := os.Args[1]
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	packageName = f.Name.Name
	fmt.Println("package " + packageName)
	fmt.Printf(`import (
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

	os.Exit(0)

	for _, s := range f.Decls {
		ast.Inspect(s, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.GenDecl:
				AnalyzeStruct(n)
			default:
				fmt.Printf("skipping %T\n", n)
			}
			fmt.Println()
			return false
		})
	}

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

	if err := templ.Execute(os.Stdout, t); err != nil {
		fmt.Printf("Error executing template: %s", err)
		os.Exit(2)
	}
}

func generateTypeOld(cmdline string) {
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
		Name  string
		Types []string
	}
	t := newtype{
		Name:  cmd[2],
		Types: strings.Split(cmd[0], "|"),
	}
	funcs := template.FuncMap{"join": strings.Join}
	templ := template.Must(template.New("class").Funcs(funcs).Parse(`
// {{.Name}} is an intersection type of {{join .Types ", "}}
type {{.Name}} struct {
{{range $i, $t := .Types}}	t{{$i}}	{{$t}}
{{end}}
	kind	int
}

func (r *{{.Name}}) String() string {
{{range $i, $t := .Types}}
	if r.kind == {{$i}} + 1 {
		return fmt.Sprintf("%v", r.t{{$i}})
	}
{{end}}
	panic("Invalid {{.Name}}")
}

// UnmarshalJSON implements json.Unmarshaler interface
func (r *{{.Name}}) UnmarshalJSON(data []byte) error {
{{range $i, $t := .Types}}
	var r{{$i}} {{$t}}
	if err := json.Unmarshal(data, &r{{$i}}); err == nil {
		r.t{{$i}} = r{{$i}}
		r.kind = {{$i}} + 1
		return nil
	}
{{end}}
	r.kind = 0
	return errors.New("invalid input: expected one of: {{join .Types ", "}}")
}

// Value returns the value of the sum-type
func (r *{{.Name}}) Value() interface{} {
{{range $i, $t := .Types}}
	if r.kind == {{$i}} + 1 {
		return r.t{{$i}}
	}
{{end}}
	panic("Invalid {{.Name}}")
}

// MarshalJSON implements json.Marshaler interface
func (r {{.Name}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Value())
}

`))

	if err := templ.Execute(os.Stdout, t); err != nil {
		fmt.Printf("Error executing template: %s", err)
		os.Exit(2)
	}
}

func AnalyzeStruct(decl *ast.GenDecl) {
	if decl.Tok != token.TYPE {
		fmt.Println("skipping " + decl.Tok.String())
		return
	}

	doc := decl.Doc
	if doc == nil {
		fmt.Println("skipping type declaration without comments")
		return
	}
	fmt.Println("analyzing type declaration")
	fmt.Printf("   doc=%s", doc.Text())
	for _, spec := range decl.Specs {
		switch sp := spec.(type) {
		case *ast.TypeSpec:
			fmt.Println("   name:", sp.Name)
			fmt.Printf("   type: %T\n", sp.Type)
		default:
			fmt.Printf("   skipping %T\n", sp)
		}
	}
}
