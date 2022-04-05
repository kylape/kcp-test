package main

import (
	"fmt"
	"strings"
)

type htmlbuilder struct {
	strings.Builder
}

func (h *htmlbuilder) tag(s string, inner func(*htmlbuilder)) {
	h.WriteString(fmt.Sprintf("<%s>\n", s))
	inner(h)
	h.WriteString(fmt.Sprintf("</%s>\n", s))
}

func (h *htmlbuilder) html(inner func(*htmlbuilder)) {
	h.tag("html", inner)
}

func (h *htmlbuilder) body(inner func(*htmlbuilder)) {
	h.tag("body", inner)
}

func (h *htmlbuilder) table(inner func(*htmlbuilder)) {
	h.tag("table", inner)
}

func (h *htmlbuilder) tr(inner func(*htmlbuilder)) {
	h.tag("tr", inner)
}
