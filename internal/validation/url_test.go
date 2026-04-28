package validation

import (
	"testing"

	"github.com/MrBorisT/url_shortener/internal/linkerr"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr error
	}{
		{
			name:    "valid url",
			url:     " http://www.example.com ",
			want:    "http://www.example.com",
			wantErr: nil,
		},
		{
			name:    "empty url",
			url:     "   ",
			want:    "",
			wantErr: linkerr.ErrURLRequired,
		},
		{
			name:    "invalid url",
			url:     "www.example.com",
			want:    "",
			wantErr: linkerr.ErrURLInvalid,
		},
		{
			name:    "invalid url scheme",
			url:     "ftp://www.example.com",
			want:    "",
			wantErr: linkerr.ErrURLInvalidScheme,
		},
		{
			name:    "missing host",
			url:     "http:///path",
			want:    "",
			wantErr: linkerr.ErrURLMissingHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.url)

			if tt.wantErr != err {
				t.Fatalf("error check: expected %v got %v", tt.wantErr, err)
			}

			if tt.want != got {
				t.Fatalf("url check: wanted %v, got %v", tt.want, got)
			}
		})
	}

}
