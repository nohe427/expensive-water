package prompt_test

import (
	"testing"

	"github.com/andreyvit/diff"
	"github.com/nohe427/expensive-water/prompt"
)

var startPrompt string = `Join the current summary with the following GitHub 
Issue Comments and provide a detailed sumary highlighting the struggles and 
comments that people made along the way. Output your summary in markdown format.
Ignore the current summary if it has not been created yet. Write your summary in
the ideal output format below.

BEGIN OUTPUT FORMAT	

# ORIGINAL ISSUE TITLE

Original issue description summary

# Issue facts

Put a bulleted list of facts about the issue here

# Notable comments from issue

Put a bulleted list of top three notable comments about the issue here

# Summarization

# Resolution of the issue if any

END OUTPUT FORMAT

BEGIN INPUTS

ORIGINAL ISSUE TITLE: {{ .OriginalIssueTitle }}

ORIGINAL ISSUE DESCRIPTION: {{ .OriginalIssueDescription }}

CURRENT SUMMARY: {{ .CurrentSummary }}

GITHUB ISSUE COMMENTS:
[ {{ range .Issues }}{{ .Body  }}, {{ end }}]`

var formattedOut string = `Join the current summary with the following GitHub 
Issue Comments and provide a detailed sumary highlighting the struggles and 
comments that people made along the way. Output your summary in markdown format.
Ignore the current summary if it has not been created yet. Write your summary in
the ideal output format below.

BEGIN OUTPUT FORMAT	

# ORIGINAL ISSUE TITLE

Original issue description summary

# Issue facts

Put a bulleted list of facts about the issue here

# Notable comments from issue

Put a bulleted list of top three notable comments about the issue here

# Summarization

# Resolution of the issue if any

END OUTPUT FORMAT

BEGIN INPUTS

ORIGINAL ISSUE TITLE: Sample Issue Title

ORIGINAL ISSUE DESCRIPTION: Sample Issue Description

CURRENT SUMMARY: Sample Current Summary

GITHUB ISSUE COMMENTS:
[ Sample Issue Body 1, Sample Issue Body 2, Sample Issue Body 3, ]`

type SampleIssue struct {
	OriginalIssueTitle       string
	OriginalIssueDescription string
	CurrentSummary           string
	Issues                   []Issue
}

type Issue struct {
	Body string
}

var sampleIssue = SampleIssue{
	OriginalIssueTitle:       "Sample Issue Title",
	OriginalIssueDescription: "Sample Issue Description",
	CurrentSummary:           "Sample Current Summary",
	Issues: []Issue{
		{Body: "Sample Issue Body 1"},
		{Body: "Sample Issue Body 2"},
		{Body: "Sample Issue Body 3"},
	},
}

func TestFormattedPrompt(t *testing.T) {
	p := prompt.Prompt{PromptText: startPrompt}
	out, _ := p.FormattedPrompt(sampleIssue)
	if out != formattedOut {
		t.Errorf("Expected something else %v", diff.LineDiff(out, formattedOut))
	}
}
