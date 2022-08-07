package utils

import (
	"strings"

	"github.com/pkg/errors"
)

func WrapErr(caller string, msg string, err error) error {
	caller = strings.ToUpper(caller)
	return errors.Wrapf(err, "[%s] %s", caller, msg)
}
