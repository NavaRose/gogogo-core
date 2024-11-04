package exception

type Http struct {
	Message    string `json:"message,omitempty"`
	Metadata   string `json:"metadata,omitempty"`
	StatusCode int    `json:"code"`
}

func (e Http) Error() string {
	return e.Message
}

func NewHttpError(message string, metadata string, statusCode int) Http {
	return Http{
		Message:    message,
		Metadata:   metadata,
		StatusCode: statusCode,
	}
}
