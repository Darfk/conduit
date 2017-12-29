package conduit

import (
	"io"
)

type stage struct {
	name   string
	input  *string
	output *string
}

func (q *stage) Name() string {
	return q.name
}

func (q *stage) IsStage() bool {
	return q.input != nil && q.output != nil
}

func (q *stage) IsSink() bool {
	return q.input == nil && q.output != nil
}

func (q *stage) IsSource() bool {
	return q.input != nil && q.output == nil
}

func (q *stage) IsDummy() bool {
	return q.input == nil && q.output == nil
}

func (q *stage) InputType() string {
	if q.input == nil {
		panic("called Input() on source")
	}
	return *q.input
}

func (q *stage) OutputType() string {
	if q.output == nil {
		panic("called Output() on sink")
	}
	return *q.output
}

func (q *stage) Execute(w io.Writer) error {
	err := template.ExecuteTemplate(w, "stage", q)
	return err
}
