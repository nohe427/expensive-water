# expensive-water

This is not an official Google product.

This is not an officially supported Google project.

## What is this?

Expensive water summarizes your GitHub issues by cataloging an issue
and all of its comments on that issue. The issue is then fed to
Google Gemini to output a summarized version of the issue and rendered
in the terminal via Markdown.

## How do I use this?

./expensive-water sum -h

Summarize a GitHub Issue and its comments. This command will make a call
to the Gemini API to retrieve the issues and all of the comments on the
issue. It will then summarize the issue and all of the comments in a neat-o
fashion.

Usage:
  expensive-water sum [flags]

Flags:
  -h, --help          help for sum
  -i, --issue int     The GitHub issue to summarize (default 1)
  -o, --org string    The GitHub org to summarize (default "firebase")
  -r, --repo string   The GitHub repo to summarize (default "flutterfire")

## Where do I put my API keys?

`${HOME}/.config/expensive-water.json`

## Can I use Vertex AI with this?

[TBD](https://github.com/nohe427/expensive-water/issues/13)

I would like to add it but its just me writing some code.
When I do get around to adding it, check the issue above for
an update.

## Why doesn't it do `X, Y, Z...`?

Feel free to add an issue to the tracker. I may not get to
it or prioritize it. This is a hobby project for me.
