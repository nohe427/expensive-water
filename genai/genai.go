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

package genai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GenClient interface {
	TokenCount(s string) int
	Summarize(s string) (string, error)
}

type GeminiClient struct {
	GenClient
	ApiKey string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{ApiKey: apiKey}
}

func (gc *GeminiClient) TokenCount(s string) int {
	ctx := context.Background()
	aic, err := genai.NewClient(ctx, option.WithAPIKey(gc.ApiKey))
	if err != nil {
		panic(err)
	}
	model := aic.GenerativeModel("gemini-pro")
	resp, err := model.CountTokens(ctx, genai.Text(s))
	if err != nil {
		panic(err)
	}
	return int(resp.TotalTokens)
}

func (gc *GeminiClient) Summarize(s string) (string, error) {
	ctx := context.Background()
	aic, err := genai.NewClient(ctx, option.WithAPIKey(gc.ApiKey))
	if err != nil {
		panic(err)
	}
	model := aic.GenerativeModel("gemini-pro")
	// model.SetMaxOutputTokens(2048)
	model.SetTemperature(0.2)
	resp, err := model.GenerateContent(ctx, genai.Text(s))
	if err != nil {
		panic(err)
	}
	if len(resp.Candidates) == 0 {
		fmt.Println("No candidates found...")
		return "", checkForFailReason(resp)
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

func checkForFailReason(resp *genai.GenerateContentResponse) error {
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
