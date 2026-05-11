package soal3

import "testing"

func TestValidateBrackets(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Nested valid - curly and angle",
			input: "{{[<>[{{}}]]}}",
			want:  true,
		},
		{
			name:  "Simple mixed valid",
			input: "[{}<>]",
			want:  true,
		},
		{
			name:  "Single closing - false",
			input: "]",
			want:  false,
		},
		{
			name:  "Closing before opening - false",
			input: "][",
			want:  false,
		},
		{
			name:  "Crossing brackets - false",
			input: "[>]",
			want:  false,
		},
		{
			name:  "Single closing angle - false",
			input: "[>",
			want:  false,
		},
		{
			name:  "Empty string",
			input: "",
			want:  false,
		},
		{
			name:  "Only openers",
			input: "<{[",
			want:  false,
		},
		{
			name:  "Only closers",
			input: "}>]",
			want:  false,
		},
		{
			name:  "Single opener",
			input: "<",
			want:  false,
		},
		{
			name:  "Single closer",
			input: ">",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateBrackets(tt.input)
			if got != tt.want {
				t.Errorf("ValidateBrackets(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
