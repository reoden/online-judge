package test

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	s := uuid.NewV4().String()
	str := uuid.NewV4().String()
	println(s)
	println(str)
	println(len(s))
}
