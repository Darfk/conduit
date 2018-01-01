package conduit

import (
	"io"
)

type stage struct {
	name   string
	input  string
	output string
}

func (q *stage) Name() string {
	return q.name
}

func (q *stage) IsStage() bool {
	return q.input != "" && q.output != ""
}

func (q *stage) IsSink() bool {
	return q.input != "" && q.output == ""
}

func (q *stage) IsSource() bool {
	return q.input == "" && q.output != ""
}

func (q *stage) IsDummy() bool {
	return q.input == "" && q.output == ""
}

func (q *stage) InputType() string {
	if q.input == "" {
		panic("called InputType() on source")
	}
	return q.input
}

func (q *stage) OutputType() string {
	if q.output == "" {
		panic("called OutputType() on sink")
	}
	return q.output
}

func (q *stage) execute(w io.Writer) error {
	var err error

	if q.IsSource() {
		err = template.ExecuteTemplate(w, "source", q)
	} else if q.IsSink() {
		err = template.ExecuteTemplate(w, "sink", q)
	} else if q.IsStage() {
		err = template.ExecuteTemplate(w, "stage", q)
	}

	return err
}
