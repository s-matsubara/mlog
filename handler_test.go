package mlog

import (
	"bytes"
	"context"
	"fmt"
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
		Format: "{{.Level}} {{.Message}} {{.File}} {{.FileName}} {{.Line}} {{.Function}}",
	})
	logger := slog.New(handler)
	logger.Info("info")
	want := "INFO info\n"
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

	fmt.Println(got)
}

func TestWithAttrs(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.With("aaa", "bbb").Info("info")

	fmt.Println(got)
}

func TestGroup(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format: "{{.Level}} {{.Message}}",
	})
	logger := slog.New(handler)
	logger.Info("info", slog.Group("group", slog.Any("aaa", "bbb")))

	fmt.Println(got)
}

func TestWithGroup(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format:      "{{.Level}} {{.Message}}",
		EnableColor: true,
	})
	logger := slog.New(handler)
	logger.WithGroup("group").Info("info", slog.Any("aaa", "bbb"))

	fmt.Println(got)
}

func TestBasicFormatter(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		EnableColor:   true,
		AttrFormatter: NewBasicFormatter(),
	})
	logger := slog.New(handler)
	logger.WithGroup("group").Info("info", slog.Any("aaa", "bbb"))

	fmt.Println(got)
}

func TestYamlFormatter(t *testing.T) {
	got := &bytes.Buffer{}
	handler := New(got, &HandlerOptions{
		Format:        "{{.Level}} {{.Message}}",
		EnableColor:   true,
		AttrFormatter: NewYamlFormatter(),
	})
	logger := slog.New(handler)
	logger.WithGroup("group").Error("error", slog.Any("aaa", "bbb"))

	fmt.Println(got)
}
