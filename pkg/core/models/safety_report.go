package models

import "time"

// ReportStatus tracks the lifecycle of a safety report.
type ReportStatus string

const (
	ReportStatusCreated  ReportStatus = "Created"
	ReportStatusRejected ReportStatus = "Rejected"
	ReportStatusResolved ReportStatus = "Resolved"
)

// ReportedContentType specifies what kind of content was reported.
type ReportedContentType string

const (
	ContentTypeMessage ReportedContentType = "Message"
	ContentTypeServer  ReportedContentType = "Server"
	ContentTypeUser    ReportedContentType = "User"
)

// Report is a safety/abuse report filed by a user.
type Report struct {
	ID                ID                  `bson:"_id,omitempty" json:"id"`
	AuthorID          ID                  `bson:"author_id" json:"author_id"`
	ContentID         ID                  `bson:"content_id" json:"content_id"`
	ContentType       ReportedContentType `bson:"content_type" json:"content_type"`
	AdditionalContext *string             `bson:"additional_context,omitempty" json:"additional_context,omitempty"`
	Status            ReportStatus        `bson:"status" json:"status"`
	Notes             *string             `bson:"notes,omitempty" json:"notes,omitempty"` // moderator notes
	CreatedAt         time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time           `bson:"updated_at" json:"updated_at"`
}
