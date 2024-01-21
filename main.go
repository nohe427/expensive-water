// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/charmbracelet/glamour"
)

func main() {
	c, err := LoadConfig("")
	if err != nil {
		fmt.Println(err)
		return
	}
	githubClient := GenerateClient(c.GetGitHubKey())
	client := NewGeminiClient(c.GetGeminiKey())
	repo, err := githubClient.GetRepo("firebase", "flutterfire")
	if err != nil {
		fmt.Println(err)
		return
	}
	issue, err := githubClient.GetIssue(repo, 10593) // summarrizing breaks with 1041
	if err != nil {
		fmt.Println(err)
		return
	}

	issueComments, err := githubClient.GetIssueComments(issue, repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	lastIndex := 0
	tokenLimit := 30_720
	toSum := ""
	currentTokenCount := 0
	defaultSumStatement := `Join the current summary with the following GitHub 
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

	Inputs:

	ORIGINAL ISSUE TITLE:%v

	ORIGINAL ISSUE DESCRIPTION:%v

	CURRENT SUMMARY:%v

	GITHUB ISSUE COMMENTS:%v`
	currSum := "No Summary Yet"
	sumSoFar := fmt.Sprintf(defaultSumStatement, issue.GetTitle(), issue.GetBody(), currSum, "")
	currentTokenCount = client.TokenCount(sumSoFar)
	for {
		for _, issue := range issueComments[lastIndex:] {
			lastIndex = lastIndex + 1
			count := client.TokenCount(issue.GetBody())
			if currentTokenCount+count > tokenLimit {
				break
			}
			toSum = toSum + issue.GetBody() + "\n\n"
			currentTokenCount = currentTokenCount + count
		}
		fmt.Println("Summarizing")
		currSum, err = client.Summarize(fmt.Sprintf(defaultSumStatement, issue.GetTitle(), issue.GetBody(), currSum, toSum))
		if err != nil {
			fmt.Println(err)
		}
		toSum = ""
		currentTokenCount = client.TokenCount(fmt.Sprintf(defaultSumStatement, issue.GetTitle(), issue.GetBody(), currSum, toSum))
		if len(issueComments) == lastIndex {
			break
		}
	}
	out, err := glamour.Render(currSum, "dracula")
	fmt.Println(out)
}
