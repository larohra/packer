package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/fatih/structtag"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

var (
	typeNames  = flag.String("type", "", "comma-separated list of type names; must be set")
	output     = flag.String("output", "", "output file name; default srcdir/<type>_hcl2.go")
	trimprefix = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of stringer:\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("hcl2-schema: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{os.Getenv("GOFILE")}
	}
	fname := args[0]
	outputPath := fname[:len(fname)-2] + "hcl2spec.go"

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("ReadFile: %+v", err)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, b, parser.ParseComments)
	if err != nil {
		fmt.Printf("ParseFile: %+v", err)
		os.Exit(1)
	}

	res := []StructDef{}

	for _, t := range types {
		for _, decl := range f.Decls {
			typeDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			typeSpec, ok := typeDecl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue
			}
			structDecl, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if typeSpec.Name.String() != t {
				continue
			}
			sd := StructDef{StructName: t}
			fields := structDecl.Fields.List
			for _, field := range fields {

				fieldType := string(b[field.Type.Pos()-1 : field.Type.End()-1])
				fieldName := fieldType[strings.Index(fieldType, ".")+1:]
				if len(field.Names) > 0 {
					fieldName = field.Names[0].Name
				}

				if !unicode.IsUpper([]rune(fieldName)[0]) {
					continue
				}
				if strings.Contains(fieldType, "func") {
					continue
				}
				fd := FieldDef{Name: fieldName}

				squash := false
				accessor := strings.ToLower(fieldName)
				if field.Tag != nil {
					tag := field.Tag.Value[1:]
					tag = tag[:len(tag)-1]
					tags, err := structtag.Parse(tag)
					if err != nil {
						log.Fatalf("structtag.Parse(%s): err: %v", field.Tag.Value, err)
					}
					if mstr, err := tags.Get("mapstructure"); err == nil {
						if len(mstr.Options) > 0 && mstr.Options[0] == "squash" {
							squash = true
						}
						if mstr.Name != "" {
							accessor = mstr.Name
						}
					}
				}

				switch fieldType {
				case "[]string":
					fd.Spec = fmt.Sprintf("%#v", &hcldec.AttrSpec{
						Name:     accessor,
						Type:     cty.List(cty.String),
						Required: false,
					})
				case "[]int":
					fd.Spec = fmt.Sprintf("%#v", &hcldec.AttrSpec{
						Name:     accessor,
						Type:     cty.List(cty.Number),
						Required: false,
					})
				case "[]byte", "string", "time.Duration":
					fd.Spec = fmt.Sprintf("%#v", &hcldec.AttrSpec{
						Name:     accessor,
						Type:     cty.String,
						Required: false,
					})
					// fd.Type = "hcl2template.Type" + strings.Title(fieldType)
				case "int", "int32", "int64", "float":
					fd.Spec = fmt.Sprintf("%#v", &hcldec.AttrSpec{
						Name:     accessor,
						Type:     cty.Number,
						Required: false,
					})
				case "bool", "config.Trilean":
					fd.Spec = fmt.Sprintf("%#v", &hcldec.AttrSpec{
						Name:     accessor,
						Type:     cty.Bool,
						Required: false,
					})
				case "[]*string", "map[*string]*string", "map[string]string", "[][]string", "TagMap":
					// TODO(azr): implement those
					continue
				case "communicator.Config":
					// this one is manually set
					continue
				case "common.PackerConfig":
					// this one is deprecated ?
					continue
				default: // nested structures
					if squash {
						sd.Squashed = append(sd.Squashed, fieldName)
						continue
					}
					sd.Nested = append(sd.Nested, NestedFieldDef{
						FieldName: fieldName,
						TypeName:  fieldType,
						Accessor:  accessor,
					})
					continue
				}

				sd.Fields = append(sd.Fields, fd)
			}
			res = append(res, sd)
		}
	}

	output := bytes.NewBuffer(nil)

	err = structDocsTemplate.Execute(output, Output{
		Package:    f.Name.String(),
		StructDefs: res,
	})
	if err != nil {
		log.Fatalf("err templating: %v", err)
	}

	formattedBytes, err := format.Source(output.Bytes())
	if err != nil {
		log.Printf("formatting err: %v", err)
		formattedBytes = output.Bytes()
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, bytes.NewBuffer(formattedBytes))
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

type Output struct {
	Package    string
	StructDefs []StructDef
}

type FieldDef struct {
	Name string
	Spec string
}
type NestedFieldDef struct {
	TypeName  string
	FieldName string
	Accessor  string
}

type StructDef struct {
	StructName string
	Fields     []FieldDef
	Nested     []NestedFieldDef
	Squashed   []string
}

var structDocsTemplate = template.Must(template.New("structDocsTemplate").
	Funcs(template.FuncMap{
		// "indent": indent,
	}).
	Parse(`// Code generated by "hcl2-schema"; DO NOT EDIT.\n

package {{ .Package }}

import (
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/zclconf/go-cty/cty"
)
{{ range .StructDefs }}
{{ $StructName := .StructName}}
func (*{{ .StructName }}) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		{{- range .Fields}}
		"{{ .Name }}": {{ .Spec }},
		{{- end }}
		{{- range .Nested}}
		"{{ .Accessor }}": &hcldec.BlockObjectSpec{TypeName: "{{ .TypeName }}", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&{{ $StructName }}{}).{{ .FieldName }}.HCL2Spec())},
		{{- end }}
	}
	{{- range .Squashed }}
	for k,v := range (&{{ $StructName }}{}).{{ . }}.HCL2Spec() {
		s[k] = v
	}
	{{- end}}
	return s
}
{{end}}`))
