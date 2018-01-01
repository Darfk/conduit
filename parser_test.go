package conduit

import (
	"testing"
)

func TestPackage(t *testing.T) {	
	const path = "testdata/generateme.go"

	var err error
	pkg, err := packageFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("package: ", pkg)
}

func TestImports(t *testing.T) {
	const path = "testdata/generateme.go"

	var err error
	imports, err := importsFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, imprt := range imports {
		t.Log("import: ", imprt)
	}
}

func TestStages(t *testing.T) {
	const path = "testdata/generateme.go"
	
	var err error
	stages, err := stagesFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, stage := range stages {
		t.Log("stage: ", stage.name)
	}

}
