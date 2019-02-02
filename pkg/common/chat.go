package common

// Request object
type Request struct {
	MessageType string   `json:"type"`
	EventTime   string   `json:"eventTime"`
	Token       string   `json:"token"`
	Message     *Message `json:"message"`
}

// Message object
type Message struct {
	Name         string `json:"name,omitempty"`
	Sender       *User  `json:"sender,omitempty"`
	CreateTime   string `json:"createTime,omitempty"`
	Text         string `json:"text,omitempty"`
	ArgumentText string `json:"argumentText,omitempty"`
}

// User object
type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
}

// Result object
type Result struct {
	Cards []*Card `json:"cards" xml:"cards"`
}

// Card object
type Card struct {
	Header   *Header    `json:"header,omitempty"`
	Sections []*Section `json:"sections"`
}

// Header object
type Header struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// Section object
type Section struct {
	Header  string    `json:"header,omitempty"`
	Widgets []*Widget `json:"widgets,omitempty"`
}

// Widget object
type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	Image         *Image         `json:"image,omitempty"`
}

// TextParagraph object
type TextParagraph struct {
	Text string `json:"text"`
}

// Image object
type Image struct {
	ImageURL string `json:"imageUrl,omitempty"`
}
