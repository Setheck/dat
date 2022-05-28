package clipper

import (
	"github.com/atotto/clipboard"
)

// Clipper is an interface for the clipboard utility
type Clipper interface {
	ReadAll() (string, error)
	WriteAll(s string) error
}

// ClipboardHelper is the main usage point
var ClipboardHelper Clipper = defaultClipper{}

type defaultClipper struct{}

// ReadAll read all data from the clipboard
func (c defaultClipper) ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write the given string to the clipboard
func (c defaultClipper) WriteAll(s string) error {
	return clipboard.WriteAll(s)
}
