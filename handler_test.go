package mlog

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestEnabled(t *testing.T) {
	ctx := context.Background()
	handler := New(os.Stdout, &HandlerOptions{
		Level: slog.LevelInfo,
	})

	var got, want bool

	got = handler.Enabled(ctx, slog.LevelDebug)
	want = false
	if got != want {
		t.Errorf("%v, want %v", got, want)
	}

	got = handler.Enabled(ctx, slog.LevelInfo)
	want = true
	if got != want {
		t.Errorf("%v, want %v", got, want)
	}

	got = handler.Enabled(ctx, slog.LevelWarn)
	want = true
	if got != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestInfo(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.Info("info")
	want := "INFO info\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestFile(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}} {{.FileName}} {{.Function}}",
	})
	logger := slog.New(handler)
	logger.Info("info")
	want := "INFO info handler_test.go github.com/s-matsubara/mlog.TestFile\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestAttrs(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.Info("info", "aaa", "bbb")
	want := "INFO info aaa=bbb\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestWithAttrs(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.With("aaa", "bbb").Info("info")
	want := "INFO info aaa=bbb\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestGroup(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.Info("info", slog.Group("group", slog.Any("aaa", "bbb")))
	want := "INFO info group=[aaa=bbb]\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestWithGroup(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.WithGroup("group").Info("info", slog.Any("aaa", "bbb"))
	want := "INFO info group=[aaa=bbb]\n"
	if got.String() != want {
		t.Errorf("%v, want %v", got, want)
	}
}

func TestColorFormatter(t *testing.T) {
	handler := New(os.Stdout, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		EnableColor:   true,
		AttrFormatter: NewBasicFormatter(),
	})
	logger := slog.New(handler)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
}

func TestBasicFormatter(t *testing.T) {
	handler := New(os.Stdout, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		AttrFormatter: NewBasicFormatter(),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slogLogger()
}

func TestNewLineFormatter(t *testing.T) {
	handler := New(os.Stdout, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		AttrFormatter: NewNewLineFormatter(),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slogLogger()
}

func TestYamlFormatter(t *testing.T) {
	handler := New(os.Stdout, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		AttrFormatter: NewYamlFormatter(),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slogLogger()
}

func slogLogger() {
	slog.Info(
		"logger",
		slog.Any("aaa", "bbb"),
		slog.Any("ccc", "ddd"),
	)
	slog.Info(
		"logger",
		slog.Any("aaa", "bbb"),
		slog.Any("ccc", "ddd"),
	)
	slog.Info(
		"logger",
		slog.Group("group",
			slog.Any("aaa", "bbb"),
			slog.Any("ccc", "ddd"),
		),
	)
	slog.Info(
		"logger",
		slog.Group("group",
			slog.Any("aaa", "bbb"),
			slog.Any("ccc", "ddd"),
		),
	)
}
