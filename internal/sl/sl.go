package sl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

type Wrapped struct {
	*slog.Logger
}

var Instance *Wrapped

func (wl *Wrapped) Fatal(msg, key string, err error) {
	wl.Error(msg, slog.String(key, err.Error()))
	os.Exit(1)
}

func init() {
	handler := NewPrettyHandler(os.Stdout)
	Instance = &Wrapped{Logger: slog.New(handler)}
}

func NewPrettyHandler(out io.Writer) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}),
		l:       log.New(out, "", 0),
	}

	return h
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	timeStr := r.Time.Format("2006/01/02 15:04:05")
	msg := color.CyanString(r.Message)

	if r.NumAttrs() > 0 {
		fields := make(map[string]any, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			fields[a.Key] = a.Value.Any()
			return true
		})

		b, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}

		h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))
	} else {
		h.l.Println(timeStr, level, msg)
	}

	return nil
}
