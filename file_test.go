package conduit

import (
	"bytes"
	"testing"
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

	buf := &bytes.Buffer{}

	createFile(buf, path, pkg, imports, stages)

	t.Log(buf.String())
}
