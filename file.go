package conduit

import (
	"bytes"
	"golang.org/x/tools/imports"
	"io"
	"os"
)

func createFile(w io.Writer, path string, pkg string, imprts []string, stages []*stage) error {

	buf := &bytes.Buffer{}

	var err error

	err = template.ExecuteTemplate(buf, "package", pkg)
	if err != nil {
		return err
	}

	err = template.ExecuteTemplate(buf, "imports", imprts)
	if err != nil {
		return err
	}

	for _, stage := range stages {
		err = stage.execute(buf)
		if err != nil {
			return err
		}
	}

	formattedBytes, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	w.Write(formattedBytes)

	return nil
}

func CreateConduitFile(dst, src string) error {
	var err error

	pkg, err := packageFromFile(src)
	if err != nil {
		return err
	}

	imprts, err := importsFromFile(src)
	if err != nil {
		return err
	}

	stages, err := stagesFromFile(src)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}

	err = createFile(buf, dst, pkg, imprts, stages)
	if err != nil {
		return err
	}

	fd, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, buf)

	return err
}
