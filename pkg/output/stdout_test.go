package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessage(t *testing.T) {
	// Given
	var buf bytes.Buffer
	var testee Outputter = SStdOut{
		Out: &buf,
	}

	// When
	testee.Error("Error message")

	// Then
	output := buf.String()

	assert.Equal(t, output, "\x1b[31mERROR\x1b[0m   Error message\n")
}

func TestWarnMessage(t *testing.T) {
	// Given
	var buf bytes.Buffer
	testee := SStdOut{
		Out: &buf,
	}

	// When
	testee.Warn("Warn message")

	// Then
	output := buf.String()

	assert.Equal(t, output, "\x1b[33mWARN\x1b[0m    Warn message\n")
}

func TestSuccessMessage(t *testing.T) {
	// Given
	var buf bytes.Buffer
	testee := SStdOut{
		Out: &buf,
	}

	// When
	testee.Success("Success message")

	// Then
	output := buf.String()

	assert.Equal(t, output, "\x1b[36mSUCCESS\x1b[0m Success message\n")
}
