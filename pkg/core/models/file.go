package models

import "time"

// FileMetadata holds type-specific metadata for an uploaded file.
type FileMetadata struct {
	// "File" | "Text" | "Image" | "Video" | "Audio"
	Type string `bson:"type" json:"type"`

	// Image / Video
	Width  *int `bson:"width,omitempty" json:"width,omitempty"`
	Height *int `bson:"height,omitempty" json:"height,omitempty"`

	// Video / Audio
	Duration *float64 `bson:"duration,omitempty" json:"duration,omitempty"` // seconds
}

// File represents an uploaded file stored in autumn (the file service).
type File struct {
	ID          ID           `bson:"_id,omitempty" json:"id"`
	Tag         string       `bson:"tag" json:"tag"` // bucket tag: "attachments", "avatars", etc.
	Filename    string       `bson:"filename" json:"filename"`
	Metadata    FileMetadata `bson:"metadata" json:"metadata"`
	ContentType string       `bson:"content_type" json:"content_type"` // MIME type
	Size        int64        `bson:"size" json:"size"`                 // bytes
	UploaderID  *ID          `bson:"uploader_id,omitempty" json:"uploader_id,omitempty"`
	Deleted     bool         `bson:"deleted" json:"deleted"`
	CreatedAt   time.Time    `bson:"created_at" json:"created_at"`
}
