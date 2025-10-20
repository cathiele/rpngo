package parse

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {
	data := []struct {
		val     string
		want    []string
		wantErr error
	}{
		{
			val: "",
		},
		{
			val:  "1",
			want: []string{"1"},
		},
		{
			val:  "\n   a bc  def\n    ",
			want: []string{"a", "bc", "def"},
		},
		{
			val:  "   \"a bc  def\"    ",
			want: []string{"\"a bc  def\""},
		},
		{
			val:  "'a bc  def'",
			want: []string{"'a bc  def'"},
		},
		{
			val:  "'a\\nb'",
			want: []string{"'a\nb'"},
		},
		{
			val:     "'a",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "a '",
			wantErr: ErrUnterminatedSingleQuote,
			want:    []string{"a"},
		},
		{
			val:     "'",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "' # comment",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "\"a",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val:     "a \"",
			wantErr: ErrUnterminatedDouble,
			want:    []string{"a"},
		},
		{
			val:     "\"",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val:     "\" # comment",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val: "# comment",
		},
		{
			val: "\n# comment \"",
		},
		{
			val:  "'\"nested quote\"'",
			want: []string{"'\"nested quote\"'"},
		},
		{
			val:  "\"'nested quote'\"",
			want: []string{"\"'nested quote'\""},
		},
		{
			val:  "\\\"",
			want: []string{"\""},
		},
		{
			val:  "\\'",
			want: []string{"'"},
		},
		{
			val:  "\\\\",
			want: []string{"\\"},
		},
		{
			val: "\\",
		},
	}

	fields := make([]string, 16)

	for _, d := range data {
		t.Run(d.val, func(t *testing.T) {
			fields = fields[:0]
			err := Fields(d.val, func(t string) error {
				fields = append(fields, t)
				return nil
			})
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err = %v, want %v", err, d.wantErr)
			}
			if len(fields) != len(d.want) {
				t.Fatalf("\n got: %+v\nwant: %+v", fields, d.want)
			}
			for i := range fields {
				if fields[i] != d.want[i] {
					t.Fatalf("\n got: %+v\nwant: %+v", fields, d.want)
				}
			}
		})
	}
}
