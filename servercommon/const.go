package servercommon

const (
	INVALID_ID uint = 0
)

const (
	HEADER_CONTENT_TYPE   string = "Content-Type"
	APPLICATION_JSON_TYPE string = "application/json"
)

const (
	PROJECT_ID_KEY     string = "project_id"
	CONFIRM_TOKEN_KEY  string = "confirm_token"
	USER_CONTEXT_KEY   string = "user"
	DEVICE_CONTEXT_KEY string = "device"
	DEVICE_ID_KEY      string = "device_id"
	CHANNEL_ID_KEY     string = "channel_id"
)

const (
	ERROR_INTERNAL               string = "server.errors.internal"
	ERROR_NAME_EMPTY             string = "server.errors.user_name_empty"
	ERROR_NAME_TOO_SHORT         string = "server.errors.user_name_too_short"
	ERROR_PASSWORD_TOO_SHORT     string = "server.errors.password_too_short"
	ERROR_EMAIL_INVALID          string = "serer.errors.email_invalid"
	ERROR_DUPLICATE_USER_NAME    string = "server.errors.duplicate_user_name"
	ERROR_DUPLICATE_EMAIL        string = "server.errors.duplicate_email"
	ERROR_NOT_AUTHORISED         string = "server.errors.authorisation_required"
	ERROR_PROJECT_NAME_TOO_SHORT string = "server.errors.project_name_too_short"
	ERROR_BAD_REQUEST            string = "server.errors.bad_request"
	ERROR_PARSE_DEVICE_TOKEN     string = "server.errors.parse_device_token"
)

type FileType int

const (
	BASE_FOLDER            string = "users_data"
	BASE_FOLDER_ENV        string = "DATCHA_BASE_FOLDER"
	PROJECTS_FOLDER        string = "projects_data"
	PROJECT_IMAGE_FILENAME string = "project_preview"
	MAX_UPLOAD_FILE_SIZE   int64  = 32000000
)

const (
	PROJECT_IMAGE FileType = 0
)

const (
	NOTIFIER_KEY string = "notifier"
)
