// Package domain defines the core business entities and interfaces for the URL
// shortening service.
package domain

import "errors"

// ErrNotFound is returned when a URL record cannot be located by its hash.
var ErrNotFound = errors.New("url not found")
