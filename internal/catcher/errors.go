package catcher

import "errors"

var (
	ErrNoSelector = errors.New("deployment has no selector")
	ErrNoPods     = errors.New("no pods found for deployment")
)
