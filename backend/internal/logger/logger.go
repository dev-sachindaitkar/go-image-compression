package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type ColorTextHandler struct {
	slog.Handler
	w io.Writer
}

func NewColorTextHandler(w io.Writer, opts *slog.HandlerOptions) *ColorTextHandler {
	return &ColorTextHandler{
		Handler: slog.NewTextHandler(w, opts),
		w:       w,
	}
}

// 1. Core interceptor that prints the colorful log lines
func (h *ColorTextHandler) Handle(ctx context.Context, r slog.Record) error {
	const (
		reset  = "\033[0m"
		blue   = "\033[34m"
		yellow = "\033[33m"
		red    = "\033[31m"
		gray   = "\033[90m"
	)

	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = gray
	case slog.LevelInfo:
		levelColor = blue
	case slog.LevelWarn:
		levelColor = yellow
	case slog.LevelError:
		levelColor = red
	default:
		levelColor = reset
	}

	timeStr := r.Time.Format("15:04:05")
	fmt.Fprintf(h.w, "%s%s%s %s%s%s %s", gray, timeStr, reset, levelColor, r.Level.String(), reset, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(h.w, " %s%s%s=%v", gray, a.Key, reset, a.Value.Any())
		return true
	})

	fmt.Fprint(h.w, "\n")
	return nil
}

// 2. REQUIRED INTERFACE METHODS (Plumbing fallbacks to satisfy slog.Handler)
func (h *ColorTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *ColorTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorTextHandler{Handler: h.Handler.WithAttrs(attrs), w: h.w}
}

func (h *ColorTextHandler) WithGroup(name string) slog.Handler {
	return &ColorTextHandler{Handler: h.Handler.WithGroup(name), w: h.w}
}
