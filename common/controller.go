package common

import ()

type Response struct {
	Data       interface{}   `json:"data"`
	Meta       *ResponseMeta `json:"meta"`
	StatusCode int           `json:"status_code"`
}

type ResponseMeta struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func PackResponse(status int, data interface{}, meta *ResponseMeta) Response {
	return Response{
		Meta:       meta,
		Data:       data,
		StatusCode: status,
	}
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	switch code {
	case DbError:
		return "DB internal error"
	case KafkaError:
		return "KAFKA internal error"
	case DuplicateError:
		return "Duplicate key error"
	default:
		return ""
	}
}