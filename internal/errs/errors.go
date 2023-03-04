package errs

import (
	"errors"
)

var (
	MissingDefaultProfileErr = errors.New("cprl.yaml is missing [default] profile")
	MissingRepositoriesErr   = errors.New("cprl.yaml is missing [repositories]")
)
