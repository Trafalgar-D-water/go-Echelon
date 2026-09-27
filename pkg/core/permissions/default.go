package permissions

const (
	DEFAULT_PERMISSION_VIEW_ONLY Mask = Mask(ViewChannel) |
		Mask(ReadMessageHistory)

	DEFAULT_PERMISSION Mask = DEFAULT_PERMISSION_VIEW_ONLY |
		Mask(SendMessage) |
		Mask(InviteOthers) |
		Mask(SendEmbeds) |
		Mask(UploadFiles) |
		Mask(Connect) |
		Mask(Speak) |
		Mask(Listen) |
		Mask(Video)

	DEFAULT_PERMISSION_SERVER Mask = DEFAULT_PERMISSION |
		Mask(React) |
		Mask(ChangeNickname) |
		Mask(ChangeAvatar)

	ALLOW_IN_TIMEOUT Mask = Mask(ViewChannel) |
		Mask(ReadMessageHistory)
)
