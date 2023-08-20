package framework

type Response struct {
	Body     any           `json:"body"`
	Code     int           `json:"code"`
	Messages []interface{} `json:"messages,omitempty"`
	Success  bool          `json:"success"`
}

type H map[string]any

func NewResponse(code int, body any, success bool, message ...interface{}) Response {
	return Response{
		Body:     body,
		Code:     code,
		Messages: message,
		Success:  success,
	}
}
