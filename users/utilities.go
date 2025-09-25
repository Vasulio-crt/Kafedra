package users

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func findKeyByValue(m map[int]int, targetValue int) (int, bool) {
	for key, value := range m {
		if value == targetValue {
			return key, true
		}
	}
	return 0, false
}

func getTokenFromHeader(r *http.Request) int {
	token, err := strconv.Atoi(r.Header.Get("token"))
	if err != nil {
		return -1
	}
	return token
}

func whatIdUser(r *http.Request) (int, bool) {
	token := getTokenFromHeader(r)
	if token == -1 {
		return 0, false
	}
	id, ok := findKeyByValue(authorized.token, token)
	if !ok {
		return 0, false
	}
	return id, true
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func errorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	message = fmt.Sprintf(`{"message": "%s"}`, message)
	w.Write([]byte(message))
}
