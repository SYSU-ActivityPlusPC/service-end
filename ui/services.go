package ui

import (
	"net/http"
	"strconv"
)

// GetHeaderMessage will read role, id and account from header
func GetHeaderMessage(r *http.Request) (int, int, string, error) {
	role, err := strconv.Atoi(r.Header.Get("X-Role"))
	if err != nil {
		return 0, 0, "", err
	}
	id, err := strconv.Atoi(r.Header.Get("X-ID"))
	if err != nil {
		return 0, 0, "", err
	}
	account := r.Header.Get("X-Account")
	return role, id, account, nil
}
