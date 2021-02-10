package slack

import (
	"github.com/slack-go/slack"
	"github.com/ztimes2/dailybugle/internal/mrkdwn"
	"github.com/ztimes2/dailybugle/internal/newspaper"
)

// Channel provides communication with a certain Slack channel.
type Channel struct {
	client    *slack.Client
	channelID string
}

// NewChannel initializes a new Channel.
func NewChannel(apiToken, channelID string) Channel {
	return Channel{
		client:    slack.New(apiToken),
		channelID: channelID,
	}
}

// Publish edits and publishes the given newspaper issue to the Slack channel.
func (c Channel) Publish(issue newspaper.Issue) error {
	blocks := []slack.Block{
		// Adds a small empty space before the very first page.
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, " ", false, false),
			nil, nil,
		),
	}

	for _, page := range issue {
		// Turns page components into Slack message blocks.

		pageBlocks := []slack.Block{
			slack.NewHeaderBlock(slack.NewTextBlockObject(
				slack.PlainTextType,
				mrkdwn.Emoji(page.HeadlineEmojiName)+" "+page.HeadlineText,
				false, false,
			)),
		}

		pageBlocks = append(pageBlocks, page.ContentElements...)

		pageBlocks = append(pageBlocks, slack.NewContextBlock(
			"",
			slack.NewTextBlockObject(
				slack.MarkdownType,
				mrkdwn.Italic("By "+page.AuthorName),
				false, false,
			),
		))

		pageBlocks = append(pageBlocks, slack.NewDividerBlock())

		blocks = append(blocks, pageBlocks...)
	}

	blocks = append(blocks,
		// Adds a small empty space after the very last page.
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, " ", false, false),
			nil, nil,
		),
	)

	if _, _, err := c.client.PostMessage(c.channelID,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionBlocks(blocks...),
	); err != nil {
		return err
	}

	return nil
}
