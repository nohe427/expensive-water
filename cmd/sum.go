// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/nohe427/expensive-water/config"
	"github.com/nohe427/expensive-water/genai"
	"github.com/nohe427/expensive-water/gh"
	"github.com/spf13/cobra"
)

type Opt struct {
	debug  bool
	vertex bool
}

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
		opt := Opt{debug: false, vertex: false}
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
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			fmt.Println(err)
		}
		opt.debug = debug
		vertex, err := cmd.Flags().GetBool("vertex")
		if err != nil {
			fmt.Println(err)
		}
		opt.vertex = vertex

		sum(repo, issue, org, opt)
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
	sumCmd.Flags().BoolP("debug", "d", false, "Whether to output files to the debug/ folder for reviewing prompts and outputs")
	sumCmd.Flags().BoolP("vertex", "v", false, "Whether to use Vertex AI instead of Gemini")
}

func sum(repo string, issue int, org string, opt Opt) {
	c, err := config.LoadConfig("")
	if err != nil {
		fmt.Println(err)
		return
	}
	ghc := gh.GenerateClient(c.GetGitHubKey())
	var client genai.GenClient
	if opt.vertex {
		opts := genai.ClientOptions{}
		opts.ApiKey = c.GetVertexKey()
		opts.Region = c.GetRegion()
		opts.ProjectId = c.GetProjectID()
		client = genai.NewVertexClient(opts)
	} else {
		client = genai.NewGeminiClient(c.GetGeminiKey())
	}
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
	// tokenLimit := 30_720
	tokenLimit := 1_000_000
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
	prevSum := ""
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
		prevSum = fmt.Sprintf(defaultSumStatement, i.GetTitle(), i.GetBody(), currSum, toSum)
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

	if opt.debug {
		out, err := writeToDebugOutput(prevSum, "input.txt")
		if err != nil {
			fmt.Println("There was an error")
		}
		fmt.Printf("input written to %v\n", out)
		out, err = writeToDebugOutput(currSum, "output.txt")
		if err != nil {
			fmt.Println("There was an error")
		}
		fmt.Printf("input written to %v\n", out)
	}
}

var currentTmp string = ""

func writeToDebugOutput(input string, filename string) (string, error) {
	if currentTmp == "" {
		tmpDir, err := os.MkdirTemp(os.TempDir(), "expensive-water")
		if err != nil {
			return "", err
		}
		currentTmp = tmpDir
	}
	outFile := filepath.Join(currentTmp, filename)
	err := os.WriteFile(outFile, []byte(input), 0644)
	if err != nil {
		return "", err
	}
	return outFile, nil
}
