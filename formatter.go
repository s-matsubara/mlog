package mlog

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"log/slog"
)

type AttrFormatter interface {
	Format(buf *bytes.Buffer, attrs []slog.Attr)
}

type basicFormatter struct{}

func NewBasicFormatter() AttrFormatter {
	return basicFormatter{}
}

func (f basicFormatter) Format(buf *bytes.Buffer, attrs []slog.Attr) {
	var b []byte

	for i, attr := range attrs {
		if attr.Equal(slog.Attr{}) {
			continue
		}
		b = append(b, []byte(attr.String())...)
		if i < len(attrs)-1 {
			b = append(b, ' ')
		}
	}

	buf.WriteByte(' ')
	buf.Write(b)
}

type yamlFormatter struct{}

func NewYamlFormatter() AttrFormatter {
	return yamlFormatter{}
}

func (f yamlFormatter) Format(buf *bytes.Buffer, attrs []slog.Attr) {
	buf.WriteByte('\n')

	f.convert(attrs)
	bs, _ := yaml.Marshal(f.convert(attrs))
	buf.Write(bs)
}

func (f yamlFormatter) convert(attrs []slog.Attr) map[string]interface{} {
	t := make(map[string]interface{})
	for _, attr := range attrs {
		switch attr.Value.Kind() {
		case slog.KindGroup:
			t[attr.Key] = f.convert(attr.Value.Group())
		default:
			t[attr.Key] = attr.Value.String()
		}
	}
	return t
}
