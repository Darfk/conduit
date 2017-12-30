package conduit

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
)

var hintRE = regexp.MustCompile(`\/\/ ?conduit`)


func packageFromFile(path string) (pkg string, err error) {
	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, path, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("could not parse file: %v\n", err)
	}

	return file.Name.Name, nil
}

func importsFromFile(path string) (imports []string, err error) {
	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return imports, fmt.Errorf("could not parse file: %v\n", err)
	}

	for _, imprt := range file.Imports {
		if imprt == nil || imprt.Path == nil {
			continue
		}
		imports = append(imports, imprt.Path.Value)
	}

	return
}

func stagesFromFile(path string) (stages []*stage, err error) {
	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return stages, fmt.Errorf("could not parse file: %v\n", err)
	}

	for _, decl := range file.Decls {
		if funcDecl, isFuncDecl := decl.(*ast.FuncDecl); isFuncDecl {
			if funcDecl.Type == nil ||
				funcDecl.Doc == nil ||
				len(funcDecl.Type.Params.List) != 1 ||
				funcDecl.Type.Results == nil ||
				len(funcDecl.Type.Results.List) != 1 {
				continue
			}

			t := funcDecl.Type

			for _, c := range funcDecl.Doc.List {
				if hintRE.MatchString(c.Text) {
					p := t.Params.List[0]
					r := t.Results.List[0]
					tmpI := &bytes.Buffer{}
					format.Node(tmpI, fs, p.Type)
					tmpO := &bytes.Buffer{}
					format.Node(tmpO, fs, r.Type)

					input := string(tmpI.Bytes())
					output := string(tmpO.Bytes())

					stages = append(stages, &stage{
						name:   funcDecl.Name.Name,
						input:  &input,
						output: &output,
					})
				}
			}
		}
	}

	return 

}
