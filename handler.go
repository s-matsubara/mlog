package mlog

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"
	"time"
)

type Handler struct {
	opts HandlerOptions
	mu   *sync.Mutex
	w    io.Writer

	attrs []slog.Attr
	group string
}

type HandlerOptions struct {
	Level         slog.Leveler
	Format        string
	TimeFormat    string
	EnableColor   bool
	AttrFormatter AttrFormatter
}

type Log struct {
	Timestamp string
	Level     string
	Message   string
	File      string
	FileName  string
	Line      int
	Function  string
}

func New(w io.Writer, opts *HandlerOptions) *Handler {
	o := getDefaultHandlerOptions()

	if opts.Level != nil {
		o.Level = opts.Level
	}
	if opts.Format != "" {
		o.Format = opts.Format
	}
	if opts.TimeFormat != "" {
		o.TimeFormat = opts.TimeFormat
	}
	if opts.EnableColor {
		o.EnableColor = opts.EnableColor
	}
	if opts.AttrFormatter != nil {
		o.AttrFormatter = opts.AttrFormatter
	}

	return &Handler{
		opts: *o,
		mu:   &sync.Mutex{},
		w:    w,
	}
}

func getDefaultHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		Level:         slog.LevelInfo,
		Format:        "{{.Timestamp}} {{.Level}} {{.Message}}",
		TimeFormat:    time.DateTime,
		EnableColor:   false,
		AttrFormatter: NewBasicFormatter(),
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.opts.Level.Level() <= level
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	h = h.clone()
	buf := &bytes.Buffer{}

	log := &Log{
		Timestamp: record.Time.Format(h.opts.TimeFormat),
		Level:     record.Level.String(),
		Message:   record.Message,
	}

	if record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		log.File = f.File
		log.FileName = filepath.Base(f.File)
		log.Line = f.Line
		log.Function = f.Function
	}

	if h.opts.EnableColor {
		log.Level = h.coloring(record)
	}

	h.mergeWithAttrs(record)

	if h.group != "" {
		h.mergeWithGroup()
	}

	tmpl, _ := template.New("log").Parse(h.opts.Format)
	if err := tmpl.Execute(buf, log); err != nil {
		return err
	}

	h.opts.AttrFormatter.Format(buf, h.attrs)

	return h.execute(buf)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h = h.clone()
	h.attrs = attrs
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h = h.clone()
	h.group = name
	return h
}

func (h *Handler) mergeWithAttrs(record slog.Record) {
	record.Attrs(func(attr slog.Attr) bool {
		h.attrs = append(h.attrs, attr)
		return true
	})
}

func (h *Handler) coloring(record slog.Record) string {
	switch record.Level {
	case slog.LevelDebug:
		return color.CyanString(h.opts.Level.Level().String())
	case slog.LevelInfo:
		return color.BlueString(h.opts.Level.Level().String())
	case slog.LevelWarn:
		return color.YellowString(h.opts.Level.Level().String())
	case slog.LevelError:
		return color.RedString(h.opts.Level.Level().String())
	}
	return ""
}

func (h *Handler) mergeWithGroup() {
	var attrs []slog.Attr
	a := slog.Attr{
		Key:   h.group,
		Value: slog.AnyValue(h.attrs),
	}
	h.attrs = append(attrs, a)
}

func (h *Handler) clone() *Handler {
	return &Handler{
		opts:  h.opts,
		mu:    h.mu,
		w:     h.w,
		attrs: h.attrs,
		group: h.group,
	}
}

func (h *Handler) execute(buf *bytes.Buffer) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
