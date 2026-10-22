package charmimage

import "testing"

func TestDetectFromMatrix(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Protocol
	}{
		{
			name: "ambiente vazio",
			env:  map[string]string{},
			want: ProtocolHalfBlock,
		},
		{
			name: "kitty pelo identificador de janela",
			env:  map[string]string{"KITTY_WINDOW_ID": "1"},
			want: ProtocolKitty,
		},
		{
			name: "iterm2 pelo TERM_PROGRAM",
			env:  map[string]string{"TERM_PROGRAM": "iTerm.app"},
			want: ProtocolITerm2,
		},
		{
			name: "wezterm pelo TERM_PROGRAM",
			env:  map[string]string{"TERM_PROGRAM": "WezTerm"},
			want: ProtocolITerm2,
		},
		{
			name: "mintty pelo TERM_PROGRAM",
			env:  map[string]string{"TERM_PROGRAM": "mintty"},
			want: ProtocolITerm2,
		},
		{
			name: "sixel no TERM",
			env:  map[string]string{"TERM": "xterm-sixel"},
			want: ProtocolSixel,
		},
		{
			name: "sixel no COLORTERM",
			env:  map[string]string{"TERM": "xterm-256color", "COLORTERM": "sixel"},
			want: ProtocolSixel,
		},
		{
			name: "terminal comum com truecolor",
			env:  map[string]string{"TERM": "xterm-256color", "COLORTERM": "truecolor"},
			want: ProtocolHalfBlock,
		},
		{
			name: "TERM dumb",
			env:  map[string]string{"TERM": "dumb"},
			want: ProtocolNone,
		},
		{
			name: "NO_COLOR definido",
			env:  map[string]string{"NO_COLOR": "1", "KITTY_WINDOW_ID": "1"},
			want: ProtocolNone,
		},
		{
			name: "dumb vence kitty",
			env:  map[string]string{"TERM": "dumb", "KITTY_WINDOW_ID": "1"},
			want: ProtocolNone,
		},
		{
			name: "kitty vence iterm2",
			env:  map[string]string{"KITTY_WINDOW_ID": "1", "TERM_PROGRAM": "iTerm.app"},
			want: ProtocolKitty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFrom(func(key string) string {
				return tt.env[key]
			})
			if got != tt.want {
				t.Fatalf("DetectFrom() = %v, quer %v", got, tt.want)
			}
		})
	}
}

func TestDetectFromNil(t *testing.T) {
	if got := DetectFrom(nil); got != ProtocolHalfBlock {
		t.Fatalf("DetectFrom(nil) = %v, quer %v", got, ProtocolHalfBlock)
	}
}

func TestDetectUsesProcessEnvironment(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("COLORTERM", "")
	t.Setenv("KITTY_WINDOW_ID", "42")
	t.Setenv("NO_COLOR", "")

	if got := Detect(); got != ProtocolKitty {
		t.Fatalf("Detect() = %v, quer %v", got, ProtocolKitty)
	}
}

func TestProtocolString(t *testing.T) {
	tests := []struct {
		protocol Protocol
		want     string
	}{
		{ProtocolNone, "none"},
		{ProtocolHalfBlock, "halfblock"},
		{ProtocolSixel, "sixel"},
		{ProtocolKitty, "kitty"},
		{ProtocolITerm2, "iterm2"},
		{Protocol(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.protocol.String(); got != tt.want {
			t.Errorf("Protocol(%d).String() = %q, quer %q", int(tt.protocol), got, tt.want)
		}
	}
}
