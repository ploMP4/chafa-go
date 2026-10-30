package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/ploMP4/chafa-go"
)

func TestPlaybackLoops(t *testing.T) {
	cursorHome, restoreTerminal := terminalCommands(t)
	for _, test := range []struct {
		name      string
		loopCount int
		once      bool
		want      int
	}{
		{"once", -1, false, 2},
		{"repeat", 2, false, 6},
		{"override", 0, true, 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			animation := &chafa.Animation{
				Width: 1, Height: 1, LoopCount: test.loopCount,
				Frames: []chafa.AnimationFrame{
					{Pixels: []uint8{255, 0, 0, 255}, Delay: time.Nanosecond},
					{Pixels: []uint8{0, 0, 255, 255}, Delay: time.Nanosecond},
				},
			}
			var output bytes.Buffer
			if err := play(context.Background(), &output, animation, 10, 5, test.once); err != nil {
				t.Fatal(err)
			}
			if got := strings.Count(output.String(), cursorHome); got != test.want {
				t.Errorf("rendered %d frames, want %d", got, test.want)
			}
			if !strings.HasSuffix(output.String(), restoreTerminal) {
				t.Error("terminal state was not restored")
			}
		})
	}
}

func TestPlaybackCancellation(t *testing.T) {
	cursorHome, restoreTerminal := terminalCommands(t)
	animation := &chafa.Animation{
		Width: 1, Height: 1,
		Frames: []chafa.AnimationFrame{{Pixels: []uint8{255, 0, 0, 255}, Delay: time.Hour}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output := &cancelWriter{cancel: cancel, cursorHome: cursorHome}
	if err := play(ctx, output, animation, 10, 5, false); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(output.String(), restoreTerminal) {
		t.Error("cancellation did not restore terminal state")
	}
}

func TestPlaybackFrameDelays(t *testing.T) {
	for _, test := range []struct {
		name  string
		delay time.Duration
		want  time.Duration
	}{
		{"untimed", 0, 50 * time.Millisecond},
		{"timed", 30 * time.Millisecond, 30 * time.Millisecond},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				animation := &chafa.Animation{
					Width: 1, Height: 1,
					Frames: []chafa.AnimationFrame{{Pixels: []uint8{255, 0, 0, 255}, Delay: test.delay}},
				}
				started := time.Now()
				var output bytes.Buffer
				if err := play(context.Background(), &output, animation, 10, 5, true); err != nil {
					t.Fatal(err)
				}
				if elapsed := time.Since(started); elapsed != test.want {
					t.Errorf("playback took %v, want %v", elapsed, test.want)
				}
			})
		})
	}
}

func terminalCommands(t *testing.T) (cursorHome, restoreTerminal string) {
	t.Helper()
	termInfo := chafa.TermDbDetect(chafa.TermDbGetDefault(), os.Environ())
	defer chafa.TermInfoUnref(termInfo)
	cursorHome = chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_CURSOR_TO_TOP_LEFT)
	if cursorHome == "" {
		t.Fatal("test terminal does not provide a cursor home command")
	}
	restoreTerminal = chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_RESET_ATTRIBUTES) +
		chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_ENABLE_CURSOR) +
		chafa.TermInfoGetSeq(termInfo, chafa.CHAFA_TERM_SEQ_DISABLE_ALT_SCREEN)
	return cursorHome, restoreTerminal
}

type cancelWriter struct {
	bytes.Buffer
	cancel     context.CancelFunc
	cursorHome string
}

func (w *cancelWriter) Write(p []byte) (int, error) {
	if bytes.HasPrefix(p, []byte(w.cursorHome)) {
		w.cancel()
	}
	return w.Buffer.Write(p)
}
