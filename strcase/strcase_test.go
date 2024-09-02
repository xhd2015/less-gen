package strcase

import (
	"reflect"
	"testing"
)

func Test_splitCamelCase(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{
			s:    "",
			want: nil,
		},
		{
			s:    "A",
			want: []string{"A"},
		},
		{
			s:    "a",
			want: []string{"a"},
		},
		{
			s:    "aDb",
			want: []string{"a", "Db"},
		},
		{
			s:    "aDBALook",
			want: []string{"a", "DBA", "Look"},
		},
		{
			s:    "DBALook",
			want: []string{"DBA", "Look"},
		},
		{
			s:    "DBALookA",
			want: []string{"DBA", "Look", "A"},
		},
		{
			s:    "HTTPSProtocol",
			want: []string{"HTTPS", "Protocol"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := SplitCamelCase(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
