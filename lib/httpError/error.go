package httpError

type HTTPError struct {
	StatusCode int `json:"statusCode"`

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}
