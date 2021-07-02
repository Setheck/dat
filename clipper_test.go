package main

import (
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClipper(t *testing.T) {
	if clipboard.Unsupported {
		t.SkipNow()
	}
	testClipData := "some clipboard data"
	if err := ClipboardHelper.WriteAll(testClipData); err != nil {
		t.Fatal(err)
	}
	if got, err := ClipboardHelper.ReadAll(); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, testClipData, got)
	}
}
