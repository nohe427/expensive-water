/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/nohe427/expensive-water/config"
	"github.com/nohe427/expensive-water/genai"
	"github.com/nohe427/expensive-water/gh"
	"github.com/spf13/cobra"
)

// sumCmd represents the sum command
var sumCmd = &cobra.Command{
	Use:   "sum",
	Short: "Summarize a GitHub Issue and its comments",
	Long: `
Summarize a GitHub Issue and its comments. This command will make a call
to the Gemini API to retrieve the issues and all of the comments on the
issue. It will then summarize the issue and all of the comments in a neat-o
fashion.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sum called")
		repo, err := cmd.Flags().GetString("repo")
		if err != nil {
			fmt.Println(err)
		}
		issue, err := cmd.Flags().GetInt("issue")
		if err != nil {
			fmt.Println(err)
		}
		org, err := cmd.Flags().GetString("org")
		if err != nil {
			fmt.Println(err)
		}
		sum(repo, issue, org)
	},
}

func init() {
	rootCmd.AddCommand(sumCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sumCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sumCmd.Flags().StringP("repo", "r", "flutterfire", "The GitHub repo to summarize")
	sumCmd.Flags().IntP("issue", "i", 1, "The GitHub issue to summarize")
	sumCmd.Flags().StringP("org", "o", "firebase", "The GitHub org to summarize")

}

func sum(repo string, issue int, org string) {
	c, err := config.LoadConfig("")
	if err != nil {
		fmt.Println(err)
		return
	}
	ghc := gh.GenerateClient(c.GetGitHubKey())
	client := genai.NewGeminiClient(c.GetGeminiKey())
	r, err := ghc.GetRepo(org, repo)
	if err != nil {
		fmt.Println(err)
		return
	}
	i, err := ghc.GetIssue(r, issue)
	if err != nil {
		fmt.Println(err)
		return
	}
	ic, err := ghc.GetIssueComments(i, r)
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
	sumSoFar := fmt.Sprintf(defaultSumStatement, i.GetTitle(), i.GetBody(), currSum, "")
	currentTokenCount = client.TokenCount(sumSoFar)
	for {
		for _, issue := range ic[lastIndex:] {
			lastIndex = lastIndex + 1
			count := client.TokenCount(issue.GetBody())
			if currentTokenCount+count > tokenLimit {
				break
			}
			toSum = toSum + issue.GetBody() + "\n\n"
			currentTokenCount = currentTokenCount + count
		}
		fmt.Println("Summarizing")
		currSum, err = client.Summarize(fmt.Sprintf(defaultSumStatement, i.GetTitle(), i.GetBody(), currSum, toSum))
		if err != nil {
			fmt.Println(err)
		}
		toSum = ""
		currentTokenCount = client.TokenCount(fmt.Sprintf(defaultSumStatement, i.GetTitle(), i.GetBody(), currSum, toSum))
		if len(ic) == lastIndex {
			break
		}
	}
	out, err := glamour.Render(currSum, "dracula")
	fmt.Println(out)
}
