package dbstorage

import "fmt"

type DatabaseError struct {
	Type string
	Err  error
}

func newDatabaseError(eType string, err error) *DatabaseError {
	return &DatabaseError{eType, err}
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("%s=%s", e.Type, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
