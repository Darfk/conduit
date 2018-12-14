package conduit

import ()

type Job interface {
	Do() []Job
	Route() int
}
