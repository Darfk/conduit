package conduit

import (
	"errors"
)

var (
	ErrPoolSizeZero = errors.New("cannot set pool size to 0")
	ErrNetworkStartedPool = errors.New("cannot set pool size after network has started")
	ErrNetworkStartedConnect = errors.New("cannot connect after network has started")
	ErrOutputConnected = errors.New("cannot connect already connected output")
)
