package slack

// Request object
type Request struct {
	Token       string `json:"token"`
	Domain      string `json:"team_domain"`
	Channel     string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Query       string `json:"text"`
	ResponseURL string `json:"response_url"`
}

// Field object
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Action object
type Action struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	URL   string `json:"url"`
	Style string `json:"style"`
}

// Attachment object
type Attachment struct {
	Fallback     *string   `json:"fallback"`
	Color        *string   `json:"color"`
	PreText      *string   `json:"pretext"`
	AuthorName   *string   `json:"author_name"`
	AuthorLink   *string   `json:"author_link"`
	AuthorIcon   *string   `json:"author_icon"`
	Title        *string   `json:"title"`
	TitleLink    *string   `json:"title_link"`
	Text         *string   `json:"text"`
	ImageURL     *string   `json:"image_url"`
	Fields       []*Field  `json:"fields"`
	Footer       *string   `json:"footer"`
	FooterIcon   *string   `json:"footer_icon"`
	Timestamp    *int64    `json:"ts"`
	MarkdownIn   *[]string `json:"mrkdwn_in"`
	Actions      []*Action `json:"actions"`
	CallbackID   *string   `json:"callback_id"`
	ThumbnailURL *string   `json:"thumb_url"`
}

// Payload object
type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}
