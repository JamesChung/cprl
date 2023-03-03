package errs

import (
	"errors"
)

var (
	ErrorCh = make(chan error, 1)

	MissingDefaultProfileErr = errors.New("cprl.yaml is missing [default] profile")
	MissingRepositoriesErr   = errors.New("cprl.yaml is missing [repositories]")
)
