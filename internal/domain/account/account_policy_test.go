package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountPolicy_LimitSavingAccount(t *testing.T) {
	policy := NewAccountPolicy()

	tests := []struct {
		name     string
		count    int
		expected bool
	}{
		{
			name:     "success - under limit",
			count:    4,
			expected: false,
		},
		{
			name:     "success - at limit",
			count:    5,
			expected: true,
		},
		{
			name:     "success - over limit",
			count:    6,
			expected: true,
		},
		{
			name:     "success - zero",
			count:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.LimitSavingAccount(tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAccountPolicy(t *testing.T) {
	policy := NewAccountPolicy()
	assert.NotNil(t, policy)
}