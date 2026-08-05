package http

import "net/http"

func userID(r *http.Request) int64 {
	if v, ok := r.Context().Value(userIDKey).(int64); ok {
		return v
	}
	return 0
}
