package models

// EmbedType classifies what kind of embed this is.
type EmbedType string

const (
	EmbedTypeWebsite EmbedType = "Website"
	EmbedTypeImage   EmbedType = "Image"
	EmbedTypeVideo   EmbedType = "Video"
	EmbedTypeGifV    EmbedType = "GifV"
	EmbedTypeNone    EmbedType = "None"
)

// ImageInfo holds width/height metadata for images inside embeds.
type ImageInfo struct {
	URL    string `bson:"url" json:"url"`
	Width  int    `bson:"width" json:"width"`
	Height int    `bson:"height" json:"height"`
}

// EmbedSpecial contains platform-specific extra data (e.g. YouTube video ID).
type EmbedSpecial struct {
	Type      string  `bson:"type" json:"type"`                               // e.g. "YouTube", "Twitch"
	ID        *string `bson:"id,omitempty" json:"id,omitempty"`               // video/channel ID
	Timestamp *int    `bson:"timestamp,omitempty" json:"timestamp,omitempty"` // seconds
}

// Embed is a rich preview attached to a message (link preview, image card, etc.).
type Embed struct {
	Type        EmbedType     `bson:"type" json:"type"`
	URL         *string       `bson:"url,omitempty" json:"url,omitempty"`
	Title       *string       `bson:"title,omitempty" json:"title,omitempty"`
	Description *string       `bson:"description,omitempty" json:"description,omitempty"`
	SiteURL     *string       `bson:"site_url,omitempty" json:"site_url,omitempty"`
	IconURL     *string       `bson:"icon_url,omitempty" json:"icon_url,omitempty"`
	Image       *ImageInfo    `bson:"image,omitempty" json:"image,omitempty"`
	Video       *ImageInfo    `bson:"video,omitempty" json:"video,omitempty"`
	Special     *EmbedSpecial `bson:"special,omitempty" json:"special,omitempty"`
	Colour      *string       `bson:"colour,omitempty" json:"colour,omitempty"` // accent hex
}
