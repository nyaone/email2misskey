package types

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
		ID      string `json:"id"`
		Kind    string `json:"kind"`
	} `json:"error"`
}

type I_Request struct {
	I string `json:"i"` // Token
}

type I_Response struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	// Ignore other fields
} // And same as type APIResponseSuccess_UserShow

type UserShow_Request struct {
	Username string  `json:"username"`
	Host     *string `json:"host"` // Null
}

// APIRequest_DriveFilesCreate is a form request

type DriveFilesCreate_Response struct {
	ID string `json:"id"`
	// Ignore other fields
}

type MessageCreate_Request struct {
	I      string `json:"i"` // Token
	UserID string `json:"userId"`
	Text   string `json:"text"`
	FileID string `json:"fileId"`
}

// No need for MessageCreate_Response
