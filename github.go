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
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
)

type GitHubClient struct {
	Client *github.Client
}

func GenerateClient(authToken string) *GitHubClient {
	ghc := &GitHubClient{}
	ghc.Client = github.NewClient(nil)

	if authToken != "" {
		fmt.Println("No auth token provided, generating default client")
		ghc.Client = ghc.Client.WithAuthToken(authToken)
	}
	return ghc
}

func (ghc *GitHubClient) GetRepo(owner string, repo string) (*github.Repository, error) {
	repository, _, err := ghc.Client.Repositories.Get(context.Background(), owner, repo)

	return repository, err
}

func (ghc *GitHubClient) GetIssue(repo *github.Repository, issueNumber int) (*github.Issue, error) {
	issue, _, err := ghc.Client.Issues.Get(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), issueNumber)

	return issue, err
}

func (ghc *GitHubClient) GetIssueComments(issue *github.Issue, repo *github.Repository) ([]*github.IssueComment, error) {
	issueComments := []*github.IssueComment{}
	opt := &github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 10}}

	for {
		issues, resp, err := ghc.Client.Issues.ListComments(
			context.Background(),
			*repo.GetOwner().Login,
			repo.GetName(),
			*issue.Number,
			opt)
		if err != nil {
			return nil, err
		}
		issueComments = append(issueComments, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return issueComments, nil
}
