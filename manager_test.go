package stm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	tm := &manager{
		tokens:   make([]token, 1),
		maxUses:  1,
		maxTime:  time.Hour,
		maxCount: 1,
	}

	token1 := tm.Get()
	assert.NotEmpty(t, token1, "Expected a non-empty token")
	assert.Equal(t, 0, tm.index, "Expected index to be zero after token1")
	assert.Equal(t, 1, tm.tokens[0].count, "Expected token count to be zero")
	token2 := tm.Get()
	assert.NotEqual(t, token1, token2, "Expected a different token")
	assert.Equal(t, 0, tm.index, "Expected index to be zero after token2")
	assert.Equal(t, 1, tm.tokens[0].count, "Expected token count to be reset")
}

func Test_GetAfterTimeout(t *testing.T) {
	tc := UUIDTokenManager(1, 2, time.Millisecond*2)

	token1 := tc.Get()
	assert.NotEmpty(t, token1, "Expected a non-empty token")
	token2 := tc.Get()
	assert.Equal(t, token1, token2, "Expected the same token")
	time.Sleep(time.Millisecond * 3)
	token3 := tc.Get()
	assert.NotEqual(t, token1, token3, "Expected a different token")
}
func Test_GetTooManyTokens(t *testing.T) {
	tc := UUIDTokenManager(3, 1, time.Second)

	token1 := tc.Get()
	assert.NotEmpty(t, token1, "Expected a non-empty token")
	token2 := tc.Get()
	assert.NotEqual(t, token1, token2, "Expected a different token")
	token3 := tc.Get()
	assert.NotEqual(t, token2, token3, "Expected a different token")
	token4 := tc.Get()
	assert.NotEqual(t, token3, token4, "Expected a different token")
}

func Test_GetTheSameTokenAgain(t *testing.T) {
	tc := UUIDTokenManager(3, 2, time.Second)

	token1 := tc.Get()
	assert.NotEmpty(t, token1, "Expected a non-empty token")
	token2 := tc.Get()
	assert.NotEqual(t, token1, token2, "Expected a different token")
	token3 := tc.Get()
	assert.NotEqual(t, token2, token3, "Expected a different token")
	token4 := tc.Get()
	assert.Equal(t, token4, token1, "Expected token4 to match token1")
}

func Test_Check(t *testing.T) {
	tm := UUIDTokenManager(1, 1, time.Hour)

	token1 := tm.Get()
	assert.True(t, tm.Check(token1), "Expected token to be valid")
}

func Test_CheckInvalid(t *testing.T) {
	tm := &manager{
		tokens:  make([]token, 1),
		maxUses:  1,
		maxTime:  time.Hour,
		maxCount: 1,
	}
	token1 := tm.Get()
	assert.True(t, tm.Check(token1), "Expected token to be valid")
	tm.tokens[0].value = UUIDToken()
	assert.False(t, tm.Check(token1), "Expected token to be invalid")
}

func Test_CheckReuse(t *testing.T) {
	tm := UUIDTokenManager(1, 2, time.Hour)
	token1 := tm.Get()
	assert.True(t, tm.Check(token1), "Expected token to be valid")
	token2 := tm.Get()
	assert.False(t, tm.Check(token1), "Expected token to be invalid")
	assert.True(t, tm.Check(token2), "Expected token to be valid")
}

func Test_CheckTooManyReuses(t *testing.T) {
	tm := UUIDTokenManager(1, 1, time.Hour)
	token1 := tm.Get()
	assert.True(t, tm.Check(token1), "Expected token to be valid")
	token2 := tm.Get()
	assert.False(t, tm.Check(token1), "Expected token to be invalid")
	assert.True(t, tm.Check(token2), "Expected token to be valid")
}
