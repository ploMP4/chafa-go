package chafa

import (
	"fmt"
	"testing"
)

func TestTermDbDetect(t *testing.T) {
	termDb := TermDbGetDefault()

	for n := 0; n <= 300; n++ {
		envp := make([]string, 0, n)
		want := "dumb"
		if n > 0 {
			envp = append(envp, "TERM=xterm-256color")
			want = "xterm-256color"
		}
		for i := len(envp); i < n; i++ {
			envp = append(envp, fmt.Sprintf("CHAFA_GO_TEST_%d=1", i))
		}

		info := TermDbDetect(termDb, envp)
		if got := TermInfoGetName(info); got != want {
			t.Fatalf("TermDbDetect() with %d variables detected %q, want %q", n, got, want)
		}
		TermInfoUnref(info)
	}
}
