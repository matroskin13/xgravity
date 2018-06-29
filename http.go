package xgravity

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Status bool
	Error  string
	Data   interface{}
}

func GetParam(url string, offset int) int {
	items := strings.Split(url, "/")

	if len(items) < offset {
		return -1
	}

	digit := items[offset]

	if v, err := strconv.Atoi(digit); err == nil {
		return v
	}

	return -1
}

func ErrorResponse(w http.ResponseWriter, message error) {
	b, _ := json.Marshal(Response{Error: message.Error()})

	w.WriteHeader(500)
	w.Write(b)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	b, _ := json.Marshal(Response{Data: data})

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
