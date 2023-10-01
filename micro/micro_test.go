package micro

import (
	"testing"
)

func TestMicroservice(t *testing.T) {
	var tests = []struct {
		name        string
		inputFile   string
		packageName string
		expected    string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
