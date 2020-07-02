package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/digitalocean/github-changelog-generator/ghcl"
)

var (
	org   = flag.String("org", "", "organization (required)")
	repo  = flag.String("repo", "", "repository (required)")
	token = flag.String("token", os.Getenv("GITHUB_TOKEN"), "Github token")
)

// FormatChangelogEntries formats a slice of changelog entries, returning a
// string that can be used to display them.
func FormatChangelogEntries(entries []*ghcl.ChangelogEntry) string {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].MergedAt.After(entries[j].MergedAt)
	})

	builder := strings.Builder{}
	for _, entry := range entries {
		builder.WriteString("- ")
		builder.WriteString(entry.Body)
		builder.WriteString(" #")
		builder.WriteString(strconv.Itoa(entry.Number))
		builder.WriteString("\n")
	}
	return builder.String()
}

func Build(cs ghcl.ChangelogService) error {
	entries, err := ghcl.FetchChangelogEntries(cs)
	if err != nil {
		return err
	}

	notes := FormatChangelogEntries(entries)
	_, err = io.WriteString(os.Stdout, notes)
	return err
}

func main() {
	flag.Parse()

	if *org == "" || *repo == "" {
		flag.Usage()
		os.Exit(1)
	}

	cs := ghcl.NewGitHubChangelogService(*org, *repo, *token, "")
	if err := Build(cs); err != nil {
		log.Fatalln(err)
	}
}
