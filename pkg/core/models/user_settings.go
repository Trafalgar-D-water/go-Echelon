package models

// UserSettings holds per-user preferences stored in the database.
type UserSettings struct {
	UserID ID `bson:"_id" json:"user_id"`

	// Appearance
	Theme    string `bson:"theme" json:"theme"`       // "dark" | "light" | "amoled"
	Language string `bson:"language" json:"language"` // e.g. "en", "fr"

	// Notifications
	NotificationsEnabled bool `bson:"notifications_enabled" json:"notifications_enabled"`
	DMNotifications      bool `bson:"dm_notifications" json:"dm_notifications"`
	MentionNotifications bool `bson:"mention_notifications" json:"mention_notifications"`

	// Privacy
	ShowOnlineStatus  bool `bson:"show_online_status" json:"show_online_status"`
	AllowDMsFromAll   bool `bson:"allow_dms_from_all" json:"allow_dms_from_all"`
	ShowCurrentServer bool `bson:"show_current_server" json:"show_current_server"`
}
