package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"log/slog"
)

type SimpleHandler struct {
	w io.Writer
}

func (h *SimpleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *SimpleHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	if r.Level >= slog.LevelError {
		b.WriteString("[")
		b.WriteString(strings.ToUpper(r.Level.String()))
		b.WriteString("] ")

		if r.PC != 0 {
			frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
			b.WriteString(frame.File)
			b.WriteString(":")
			b.WriteString(strconv.Itoa(frame.Line))
			b.WriteString(" ")
		}

		b.WriteString(r.Message)

		r.Attrs(func(attr slog.Attr) bool {
			b.WriteString(" ")
			b.WriteString(attr.Key)
			b.WriteString("=")
			b.WriteString(fmt.Sprint(attr.Value))
			return true
		})
	} else if r.Level == slog.LevelInfo {
		b.WriteString("[")
		b.WriteString(strings.ToUpper(r.Level.String()))
		b.WriteString("]: ")
		b.WriteString(r.Message)
	} else {
		b.WriteString(r.Message)
	}

	b.WriteRune('\n')
	_, err := h.w.Write([]byte(b.String()))
	return err
}

func (h *SimpleHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *SimpleHandler) WithGroup(string) slog.Handler      { return h }

func Init() {
	logger := slog.New(&SimpleHandler{w: os.Stdout})
	slog.SetDefault(logger)

	slog.Info("server started")
	slog.Error("failed to connect", "err", fmt.Errorf("timeout"))
}
