package misskey

type I_Request struct {
	I string `json:"i"` // Token
}

type I_Response struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	// Ignore other fields
} // And same as type APIResponseSuccess_UserShow

func I(accessToken string) (*I_Response, error) {
	return PostAPIRequest[I_Response]("i", &I_Request{
		I: accessToken,
	})
}
