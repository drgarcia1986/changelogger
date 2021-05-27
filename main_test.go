package main

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"time"
)

func TestWriteHeader(t *testing.T) {
	version := "v0.0.1"
	expectedVersionLine := fmt.Sprintf(
		"## %s (%s)",
		version,
		time.Now().Format("2006-01-02"),
	)

	var builder strings.Builder
	writeHeader(version, &builder)

	header := builder.String()
	for _, expected := range []string{PLACEHOLDER, expectedVersionLine} {
		if !strings.Contains(header, expected) {
			t.Errorf("want %s, got %s", expected, header)
		}
	}
}

func TestGetEntries(t *testing.T) {
	removed := []string{}
	removeFile = func(path string) error {
		removed = append(removed, path)
		return nil
	}
	expectedRemovedFiled := 5

	entries, err := getEntries("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	if actual := len(entries["added"]); actual != 2 {
		t.Errorf("want 2; got %d", actual)
	}

	expected := "The deprecated feature from the version X\n"
	if actual := entries["removed"][0]; actual != expected {
		t.Errorf("want %s; got %s", expected, actual)
	}

	if actual := len(removed); actual != expectedRemovedFiled {
		t.Errorf("want %d, got %d", expectedRemovedFiled, actual)
	}
}

func TestBuildReleaseNotes(t *testing.T) {
	removeFile = func(path string) error { return nil }
	version := "v0.0.1"

	notes, err := buildReleaseNotes(version, "./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range []string{version, "### Added", "### Fixed", "### Removed", "### Changed"} {
		if !strings.Contains(notes, expected) {
			t.Errorf("want %s, got %s", expected, notes)
		}
	}
}

func TestUpdateChangelog(t *testing.T) {
	readFile = func(_ string) ([]byte, error) {
		content := `
# Changelog

[changelogger-notes]::

## v0.0.1 (2020-05-27)
* First Version :tada:
`
		return []byte(content), nil
	}

	actualChangelog := ""
	writeFile = func(_ string, content []byte, _ fs.FileMode) error {
		actualChangelog = string(content)
		return nil
	}

	expectedChangelog := `
# Changelog


[changelogger-notes]::

## v0.0.2 (2020-05-27)
### Fixed
* A bug on some cool feature


## v0.0.1 (2020-05-27)
* First Version :tada:
`

	notes := `
[changelogger-notes]::

## v0.0.2 (2020-05-27)
### Fixed
* A bug on some cool feature
`
	if err := updateChangelog("some-path", notes); err != nil {
		t.Fatal(err)
	}

	if actualChangelog != expectedChangelog {
		t.Errorf("want %s, got %s", expectedChangelog, actualChangelog)
	}
}
