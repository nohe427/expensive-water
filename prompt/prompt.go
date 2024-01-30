package prompt

import (
	"bytes"
	"os"
	"text/template"
)

type Prompt struct {
	PromptText string
}

func LoadPrompt(location string) (*Prompt, error) {
	file, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	return &Prompt{PromptText: string(file)}, nil
}

func (p *Prompt) GetPrompt() string {
	return p.PromptText
}

func (p *Prompt) FormattedPrompt(ds any) (string, error) {
	tmpl, err := template.New("formedPrompt").Parse(p.PromptText)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	tmpl.Execute(buf, ds)
	return buf.String(), nil
}
