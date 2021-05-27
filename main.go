package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const PLACEHOLDER = "[changelogger-notes]::"

func writeHeader(version string, notesBuilder *strings.Builder) {
	fmt.Fprintf(notesBuilder, "%s\n\n", PLACEHOLDER)
	fmt.Fprintf(notesBuilder, "## %s (%s)\n", version, time.Now().Format("2006-01-02"))
}

func getEntries(entriesDir string) (map[string][]string, error) {
	entries := map[string][]string{
		"added":   []string{},
		"changed": []string{},
		"removed": []string{},
		"fixed":   []string{},
	}

	err := filepath.Walk(entriesDir, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if fi.IsDir() {
			return nil
		}

		entryType := filepath.Ext(fi.Name())[1:]
		entryList, ok := entries[entryType]
		if !ok {
			return nil
		}
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		entries[entryType] = append(entryList, string(fileContent))

		return os.Remove(path)
	})

	if err != nil {
		return nil, err
	}
	return entries, nil
}

func buildReleaseNotes(version, entriesDir string) (string, error) {
	var notesBuilder strings.Builder
	writeHeader(version, &notesBuilder)

	entries, err := getEntries(entriesDir)
	if err != nil {
		return "", err
	}

	for entryType, notes := range entries {
		if len(notes) == 0 {
			continue
		}
		fmt.Fprintf(&notesBuilder, "### %s\n", strings.Title(entryType))
		for _, note := range notes {
			fmt.Fprintf(&notesBuilder, "* %s", note)
		}
		fmt.Fprint(&notesBuilder, "\n")
	}

	return notesBuilder.String(), nil
}

func updateChangelog(path, releaseNotes string) error {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	newContent := strings.Replace(string(fileContent), PLACEHOLDER, releaseNotes, 1)

	return ioutil.WriteFile(path, []byte(newContent), 0)
}

func main() {
	version := flag.String("version", "", "The release version")
	entriesDir := flag.String("dir", "", "Directory of changelog entries")
	changelogPath := flag.String("path", "", "Path of the changelog file")

	flag.Parse()

	if *version == "" || *entriesDir == "" || *changelogPath == "" {
		flag.Usage()
		return
	}

	releaseNotes, err := buildReleaseNotes(*version, *entriesDir)
	if err != nil {
		panic(err)
	}
	if err := updateChangelog(*changelogPath, releaseNotes); err != nil {
		panic(err)
	}
}
