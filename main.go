package main

import (
	"fmt"
)

func main() {
	c, err := LoadConfig("")
	if err != nil {
		fmt.Println(err)
		return
	}
	githubClient := GenerateClient(c.Config.GitHubKey)
	client := NewGeminiClient(c.Config.GeminiKey)
	repo, err := githubClient.GetRepo("firebase", "flutterfire")
	if err != nil {
		fmt.Println(err)
		return
	}
	issue, err := githubClient.GetIssue(repo, 1041) //10593) // summarrizing breaks with 1041
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
	 
	Put a bulleted list of notable comments about the issue here

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
			fmt.Println(currentTokenCount + count)
			fmt.Println(lastIndex)
			if currentTokenCount+count > tokenLimit {
				break
			}
			toSum = toSum + issue.GetBody() + "\n\n"
			currentTokenCount = currentTokenCount + count
		}
		fmt.Println("Summarizing")
		currSum = client.Summarize(fmt.Sprintf(defaultSumStatement, issue.GetTitle(), issue.GetBody(), currSum, toSum))
		toSum = ""
		currentTokenCount = client.TokenCount(fmt.Sprintf(defaultSumStatement, issue.GetTitle(), issue.GetBody(), currSum, toSum))
		fmt.Printf("current token count %v last index %v\n", len(issueComments), lastIndex)
		if len(issueComments) == lastIndex {
			fmt.Println("breaking")
			break
		}
	}
	fmt.Println(currSum)
}
