package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringIsFlatString(t *testing.T) {
	assert := assert.New(t)
	s := NewString("foo")
	assert.IsType(s, flatString{})
}

func TestFlatStringEquals(t *testing.T) {
	assert := assert.New(t)
	s1 := NewString("foo")
	s2 := NewString("foo")
	s3 := s2
	s4 := NewString("bar")
	assert.True(s1.Equals(s2))
	assert.True(s2.Equals(s1))
	assert.True(s1.Equals(s3))
	assert.True(s3.Equals(s1))
	assert.False(s1.Equals(s4))
	assert.False(s4.Equals(s1))
}

func TestFlatStringString(t *testing.T) {
	assert := assert.New(t)
	s1 := NewString("")
	s2 := NewString("foo")
	assert.Equal("", s1.String())
	assert.Equal("foo", s2.String())
}
