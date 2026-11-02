package testing

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        2,
			b:        3,
			expected: 5,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: -5,
		},
		{
			name:     "mixed numbers",
			a:        -2,
			b:        5,
			expected: 3,
		},
		{
			name:     "zero",
			a:        0,
			b:        5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)

			if result != tt.expected {
				t.Errorf(
					"Add(%d, %d) = %d; expected %d",
					tt.a,
					tt.b,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name          string
		a             int
		b             int
		expected      int
		expectedValid bool
	}{
		{
			name:          "normal division",
			a:             10,
			b:             2,
			expected:      5,
			expectedValid: true,
		},
		{
			name:          "division with negative number",
			a:             -10,
			b:             2,
			expected:      -5,
			expectedValid: true,
		},
		{
			name:          "division by zero",
			a:             10,
			b:             0,
			expected:      0,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := Divide(tt.a, tt.b)

			if result != tt.expected {
				t.Errorf(
					"Divide(%d, %d) result = %d; expected %d",
					tt.a,
					tt.b,
					result,
					tt.expected,
				)
			}

			if valid != tt.expectedValid {
				t.Errorf(
					"Divide(%d, %d) valid = %v; expected %v",
					tt.a,
					tt.b,
					valid,
					tt.expectedValid,
				)
			}
		})
	}
}
