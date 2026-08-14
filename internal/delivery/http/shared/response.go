package http

type Response struct {
	Status  string      `json:"status"`
	Error  string      `json:"error,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func Success() Response {
	return Response{
		Status: StatusSuccess,
	}
}