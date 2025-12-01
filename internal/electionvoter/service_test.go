package electionvoter

import "testing"

func TestNormalizeAcademicStatus(t *testing.T) {
	tests := []struct {
		name      string
		input     *string
		want      string
		wantError bool
	}{
		{name: "nil defaults to active", input: nil, want: defaultAcademicStatus},
		{name: "trim and uppercase", input: strPtr(" graduated "), want: "GRADUATED"},
		{name: "empty falls back", input: strPtr(""), want: defaultAcademicStatus},
		{name: "invalid value", input: strPtr("unknown"), wantError: true},
	}

	for _, tt := range tests {
		got, err := normalizeAcademicStatus(tt.input)
		if tt.wantError {
			if err == nil {
				t.Fatalf("%s: expected error, got none", tt.name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.name, err)
		}
		if got != tt.want {
			t.Fatalf("%s: expected %s, got %s", tt.name, tt.want, got)
		}
	}
}

func strPtr(v string) *string {
	return &v
}
