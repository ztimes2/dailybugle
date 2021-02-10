package mrkdwn

// Bold formats the given string to appear bold according to Slack's mrkdwn format.
func Bold(s string) string {
	return "*" + s + "*"
}

// Italic formats the given string to appear italic according to Slack's mrkdwn
// format.
func Italic(s string) string {
	return "_" + s + "_"
}

// Link formats the given string to appear as a link to the given URL according
// to Slack's mrkdwn format.
func Link(s, url string) string {
	return "<" + url + "|" + s + ">"
}

// Emoji formats the given emoji name to an actual emoji according to Slack's
// mrkdown format.
func Emoji(name string) string {
	return ":" + name + ":"
}
