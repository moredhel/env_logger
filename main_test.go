package env_logger

import (
	"testing"
)

// TODO: Have an internal `reset` function for testing.

func TestDummy(t *testing.T) {
	ConfigureDefaultLogger()
	t.Fail()
}

