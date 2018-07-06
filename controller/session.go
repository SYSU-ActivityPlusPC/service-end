package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
	"github.com/sysu-activitypluspc/service-end/services"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Read header
	auth := r.Header.Get("Authorization")
	role := r.Header.Get("X-Role")
	if len(role) > 0 {
		IntRole, _ := strconv.Atoi(role)
		username := r.Header.Get("X-Account")
		// Handle role 1 and 2
		if IntRole == 1 || IntRole == 2 {
			user := model.GetUserInfo(username)
			res := types.PCUserInfo{
				ID:    user.ID,
				Name:  user.Name,
				Logo:  user.Logo,
				Token: auth,
			}
			stringRes, _ := json.Marshal(res)
			w.Header().Set("X-Role", role)
			w.Write(stringRes)
			return
		}
	}

	// Check user account and password
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Get post body
	jsonBody := new(types.PCUserRequest)
	err = json.Unmarshal(body, &jsonBody)
	user := model.GetUserInfo(jsonBody.Account)
	if user.ID < 0 {
		w.WriteHeader(400)
		return
	}
	role = "1"
	// Check if the user is admin
	isAdmin := services.CheckIsAdmin(jsonBody.Account)
	if isAdmin {
		role = "2"
	}
	// Validate password
	stringID := strconv.Itoa(user.ID)
	password := services.GetPassword(stringID, jsonBody.Password)
	if password == user.Password {
		// Generate token
		token, _ := services.GenerateJWT(user.Account)
		res := types.PCUserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Logo:  user.Logo,
			Token: token,
		}
		stringRes, _ := json.Marshal(res)
		w.Header().Set("X-Role", role)
		w.Write(stringRes)
		return
	}
	w.WriteHeader(400)
}