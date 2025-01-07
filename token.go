package stm

import "time"

// token is a string value, a (use) count, and time of expiration.
type token struct {
	value   string
	count   int
	expires time.Time
}
