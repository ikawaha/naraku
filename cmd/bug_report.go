package cmd

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// newBugReportCmd return bug-report command.
func newBugReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "bug-report",
		Short:   "Submit a bug report at GitHub",
		Long:    "bug-report opens the default browser to start a bug report which will include useful system information.",
		Example: "   naraku bug-report",
		RunE:    bugReport,
	}
}

func bugReport(cmd *cobra.Command, _ []string) error {
	var buf bytes.Buffer

	const (
		description = `## Description (About the problem)
A clear description of the bug encountered.

`
		toReproduce = `## Steps to reproduce
Steps to reproduce the bug.

`
		expectedBehavior = `## Expected behavior
Expected behavior.

`
		additionalDetails = `## Additional details**
Any other useful data to share.
`
	)
	buf.WriteString(fmt.Sprintf("## naraku version\n%s\n\n", cmd.Version))
	buf.WriteString(description)
	buf.WriteString(toReproduce)
	buf.WriteString(expectedBehavior)
	buf.WriteString(additionalDetails)

	body := buf.String()
	url := "https://github.com/go-spect/naraku/issues/new?title=[Bug Report] Title&body=" + url.QueryEscape(body)

	if !openBrowser(url) {
		fmt.Print("Please file a new issue at https://github.com/go-spect/naraku/issues/new using this template:\n\n")
		fmt.Print(body)
	}

	return nil
}
