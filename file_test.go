package conduit

import (
	//"bytes"
	"testing"
	"os"
)

func TestFile(t *testing.T) {
	const path = "testdata/generateme.go"

	var err error
	pkg, err := packageFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	imports, err := importsFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	stages, err := stagesFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	createFile(os.Stdout, path, pkg, imports, stages)
}
