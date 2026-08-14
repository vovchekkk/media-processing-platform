package shared

type Response struct {
	ResponseStatus  string      `json:"status"`
	Error  string      `json:"error,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Error(msg string) Response {
	return Response{
		ResponseStatus: StatusError,
		Error:  msg,
	}
}

func Success() Response {
	return Response{
		ResponseStatus: StatusSuccess,
	}
}