package shortcode

import (
	"strings"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr error
	}{
		{
			name:    "valid",
			length:  8,
			wantErr: nil,
		},
		{
			name:    "negative length",
			length:  -8,
			wantErr: ErrInvalidLength,
		},
		{
			name:    "zero length",
			length:  0,
			wantErr: ErrInvalidLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateShortCode(tt.length)

			if tt.wantErr != err {
				t.Fatalf("error check: wanted %v - got %v", tt.wantErr, err)
			}

			for _, c := range got {
				shortCodeChar := string(c)
				if !strings.Contains(alphabet, shortCodeChar) {
					t.Fatalf("character %v not in alphabet", shortCodeChar)
				}
			}
		})
	}
}
