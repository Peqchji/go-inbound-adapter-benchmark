package domain

import "fmt"

type DomainErr struct {
	DomainName string
	Err        error
}

func (e *DomainErr) Error() string {
	return fmt.Sprintf("%s: %s", e.DomainName, e.Err.Error())
}

func (e *DomainErr) Unwrap() error {
	return e.Err
}
