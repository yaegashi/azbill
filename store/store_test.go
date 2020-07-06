package store

import (
	"fmt"
	"os"
	"testing"
)

func TestLocation(t *testing.T) {
	tests := []struct {
		r          bool
		h, d, i, o string
	}{
		{r: true, h: "/home/user", d: "~/.config", i: "token.json", o: "/home/user/.config/token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "/token.json", o: "/token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "./token.json", o: "./token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "file:token.json", o: "file:token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "file:/token.json", o: "/token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "file://token.json", o: "file://token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "file:///token.json", o: "/token.json"},
		{r: true, h: "/home/user", d: "~/.config", i: "file:///./token.json", o: "/./token.json"},
		{
			r: true,
			h: "/home/user",
			d: "~/.config",
			i: "https://storage.blob.core.windows.net/container/token.json?key=secret",
			o: "https://storage.blob.core.windows.net/container/token.json?...",
		},
		{
			r: false,
			h: "/home/user",
			d: "~/.config",
			i: "https://storage.blob.core.windows.net/container/token.json?key=secret",
			o: "https://storage.blob.core.windows.net/container/token.json?key=secret",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i+1), func(t *testing.T) {
			os.Setenv("HOME", tt.h)
			s, err := NewStore(tt.d)
			if err != nil {
				t.Fatal(err)
			}
			o := s.Location(tt.i, tt.r)
			if tt.o != o {
				t.Errorf("Mismatch want %q got %q", tt.o, o)
			}
		})
	}
}
