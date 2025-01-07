package stm

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// UUIDToken returns a new UUID token.
func UUIDToken() string {
	return uuid.New().String()
}

// UUIDTokenManager creates a new TokenManager that issues UUID tokens from a ring buffer.
// maxCount is the number of tokens to allocate before the newest token replaces the oldest.
// maxUses is the maximum number of times each token can be used before it expires.
// The timeout is the aount of time before each token expires.
// Examples:
// Issue 1 token that can be used an unlimited number of times forever.
// UUIDTokenManager(1, 0, 0)
// Issue 3 tokens that can be used an unlimited number of times for up to 8 hours.
// UUIDTokenManager(3, 0, 8*time.Hour)
// Issue 5 tokens that can be used 100 times each for up to 3 hours.
// UUIDTokenManager(5, 100, 3*time.Hour)
// Issue 10 tokens that can be used fifteen times each for up to a minute.
// UUIDTokenManager(10, 15, time.Minute)
func UUIDTokenManager(maxCount, maxUses int, timeout time.Duration) TokenManager {
	tokens := make([]token, maxCount)
	for i := range tokens {
		tokens[i] = token{
			value:   UUIDToken(),
			count:   0,
			expires: time.Now().Add(timeout),
		}
	}
	if maxUses <= 0 {
		maxUses = math.MaxInt // "Unlimited"
	}
	return &manager{
		tokens:   tokens,
		maxCount: maxCount,
		maxUses:  maxUses,
		maxTime:  timeout,
		index:    0,
	}
}
