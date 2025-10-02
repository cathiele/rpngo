package parse

import (
	"errors"
	"reflect"
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

	for _, d := range data {
		t.Run(d.val, func(t *testing.T) {
			got, err := Fields(d.val)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err = %v, want %v", err, d.wantErr)
			}
			if !reflect.DeepEqual(got, d.want) {
				t.Errorf("\n got: %+v\nwant: %+v", got, d.want)
			}
		})
	}
}
