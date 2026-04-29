package domain

import "errors"

// ErrNotFound is returned when a URL record cannot be located by its hash.
var ErrNotFound = errors.New("url not found")
