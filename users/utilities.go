package users

import (
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

func getTokenFromHeader(w http.ResponseWriter, r *http.Request) int {
	token, err := strconv.Atoi(r.Header.Get("token"))
	if err != nil {
		http.Error(w, `{"error": "Invalid token"}`, http.StatusBadRequest)
		return -1
	}
	return token
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
