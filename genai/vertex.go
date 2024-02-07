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

package genai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

type VertexClient struct {
	GenClient
	Opts ClientOptions
}

type ClientOptions struct {
	ApiKey    string
	Region    string
	ProjectId string
}

var MODEL_NAME = "gemini-pro"

func NewVertexClient(opts ClientOptions) *VertexClient {
	return &VertexClient{Opts: opts}
}

func (gc *VertexClient) getClient() (*genai.Client, error) {
	ctx := context.Background()
	c, err := genai.NewClient(ctx, gc.Opts.ProjectId, gc.Opts.Region)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (gc *VertexClient) TokenCount(s string) int {
	aic, err := gc.getClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	model := aic.GenerativeModel(MODEL_NAME)
	resp, err := model.CountTokens(ctx, genai.Text(s))
	if err != nil {
		panic(err)
	}
	return int(resp.TotalTokens)
}

func (gc *VertexClient) Summarize(s string) (string, error) {
	ctx := context.Background()
	aic, err := gc.getClient()
	if err != nil {
		panic(err)
	}
	model := aic.GenerativeModel(MODEL_NAME)
	// model.SetMaxOutputTokens(2048)
	model.SetTemperature(0.2)
	resp, err := model.GenerateContent(ctx, genai.Text(s))
	if err != nil {
		panic(err)
	}
	if len(resp.Candidates) == 0 {
		fmt.Println("No candidates found...")
		return "", checkForVertexFailReason(resp)
	}
	if resp.Candidates[0].FinishReason.String() != "FinishReasonStop" {
		return "", errors.New(fmt.Sprintf("Finish reason: %v", resp.Candidates[0].FinishReason.String()))
	}
	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		v, ok := part.(genai.Text)
		if ok {
			sb.WriteString(string(v))
		}
	}
	return sb.String(), nil
}

func checkForVertexFailReason(resp *genai.GenerateContentResponse) error {
	if resp.PromptFeedback.BlockReason.String() != "" {
		return errors.New(resp.PromptFeedback.BlockReason.String())
	}
	if resp.PromptFeedback.BlockReason == 2 {
		return errors.New("prompt blocked for unknown reasons")
	}
	var sb strings.Builder
	for _, safetyRating := range resp.PromptFeedback.SafetyRatings {
		if safetyRating.Blocked {
			sb.WriteString(fmt.Sprintf("Category: %v Probability :%v\n", safetyRating.Category, safetyRating.Probability))
		}
	}
	return errors.New(sb.String())
}
