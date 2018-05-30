package controller

import (
	"sort"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/sysu-activitypluspc/service-end/model"
	"github.com/sysu-activitypluspc/service-end/types"
)

// AddActivityHandler add activity to the db
func AddActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	stringID := r.Header.Get("X-ID")
	// Handle anonymous user
	if len(stringID) == 0 {
		stringID = "-1"
	}
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.ActivityInfo)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}
	_, err = model.AddActivity(*jsonBody, ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// ModifyActivityHandler change the message except id and verified
func ModifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	if err != nil || intActID <= 0 {
		w.WriteHeader(400)
		return
	}

	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.ActivityInfo)
	json.Unmarshal(body, jsonBody)
	_, err = model.UpdateActivity(intActID, *jsonBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// DeleteActivityHandler remove activity with given id
func DeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	actid := mux.Vars(r)["actId"]
	if len(actid) <= 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	if err != nil || intActID <= 0 {
		w.WriteHeader(400)
		return
	}

	_, err = model.DeleteActivity(intActID)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

// VerifyActivityHandler change the verify status of an activity
func VerifyActivityHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actid := r.FormValue("act")
	verified := r.FormValue("verify")
	if len(actid)*len(verified) == 0 {
		w.WriteHeader(400)
		return
	}
	intActID, err := strconv.Atoi(actid)
	intVerify, err := strconv.Atoi(verified)
	if err != nil || intActID <= 0 || (intVerify != 0 && intVerify != 1) {
		w.WriteHeader(400)
		return
	}

	_, err = model.VerifyActivity(intActID, intVerify)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

// ShowActivitiesListHandler display activity with given page number and verify status
func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	// Get required page number, if not given, use the default value 1
	r.ParseForm()
	var pageNumber string
	if len(r.Form["page"]) > 0 {
		pageNumber = r.Form["page"][0]
	} else {
		pageNumber = "1"
	}
	var verified string
	if len(r.Form["verify"]) > 0 {
		verified = r.Form["verify"][0]
	} else {
		w.WriteHeader(400)
		return
	}
	intPageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}
	intVerified, err := strconv.Atoi(verified)
	if err != nil || (intVerified != 0 && intVerified != 1) {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intPageNum > 0 {
		// Get activity list
		activityList := model.GetActivityList(intPageNum-1, intVerified)

		// Change each element to the format that we need
		infoArr := make([]types.ActivityIntroduction, 0)
		for i := 0; i < len(activityList); i++ {
			tmp := types.ActivityIntroduction{
				ID:              activityList[i].ID,
				Name:            activityList[i].Name,
				StartTime:       activityList[i].StartTime.UnixNano() / int64(time.Millisecond),
				EndTime:         activityList[i].EndTime.UnixNano() / int64(time.Millisecond),
				Campus:          activityList[i].Campus,
				EnrollCondition: activityList[i].EnrollCondition,
				Sponsor:         activityList[i].Sponsor,
				PubStartTime:    activityList[i].PubStartTime.UnixNano() / int64(time.Millisecond),
				PubEndTime:      activityList[i].PubEndTime.UnixNano() / int64(time.Millisecond),
				Verified:        activityList[i].Verified,
				Type:            activityList[i].Type,
			}
			infoArr = append(infoArr, tmp)
		}
		returnList := types.ActivityList{
			Content: infoArr,
		}

		// Transfer it to json
		stringList, err := json.Marshal(returnList)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			w.WriteHeader(500)
			return
		}
		if len(activityList) <= 0 {
			w.WriteHeader(204)
		} else {
			w.Write(stringList)
		}
	} else {
		w.WriteHeader(400)
	}
}

// ShowActivityDetailHandler return required activity details with given activity id
func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		w.WriteHeader(400)
		return
	}

	// Judge if the passed param is valid
	if intID > 0 {
		ok, activityInfo := model.GetActivityInfo(intID)
		if ok {
			layout := "2006-01-02 15:04"
			// Convert to ms
			retMsg := types.ActivityInfo{
				ID:              activityInfo.ID,
				Name:            activityInfo.Name,
				StartTime:       activityInfo.StartTime.Format(layout),
				EndTime:         activityInfo.EndTime.Format(layout),
				Campus:          activityInfo.Campus,
				Location:        activityInfo.Location,
				EnrollCondition: activityInfo.EnrollCondition,
				Sponsor:         activityInfo.Sponsor,
				Type:            activityInfo.Type,
				PubStartTime:    activityInfo.PubStartTime.Format(layout),
				PubEndTime:      activityInfo.PubEndTime.Format(layout),
				Detail:          activityInfo.Detail,
				Reward:          activityInfo.Reward,
				Introduction:    activityInfo.Introduction,
				Requirement:     activityInfo.Requirement,
				Poster:          activityInfo.Poster,
				Qrcode:          activityInfo.Qrcode,
				Email:           activityInfo.Email,
				Verified:        activityInfo.Verified,
			}
			retMsg.Poster = GetPoster(retMsg.Poster, retMsg.Type)
			stringInfo, err := json.Marshal(retMsg)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				w.WriteHeader(500)
				return
			}
			w.Write(stringInfo)
		} else {
			w.WriteHeader(204)
		}
	} else {
		w.WriteHeader(400)
	}
}

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
	isAdmin := CheckIsAdmin(jsonBody.Account)
	if isAdmin {
		role = "2"
	}
	// Validate password
	stringID := strconv.Itoa(user.ID)
	password := getPassword(stringID, jsonBody.Password)
	if password == user.Password {
		// Generate token
		token, _ := GenerateJWT(user.Account)
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	jsonBody := new(types.PCUserSignInfo)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		w.WriteHeader(400)
	}
	if isEmailExist := CheckIfEmailExist(jsonBody.Email); isEmailExist == true {
		w.WriteHeader(400)
	}
	ok := model.AddUser(*jsonBody)
	if ok {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	var maxMemory int64 = 5 * (1 << 20)
	r.ParseMultipartForm(maxMemory)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
	}
	defer file.Close()
	staticFilePosition := os.Getenv("STATIC_DIR")
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	md5Filename := GetMd5(content)
	ext := path.Ext(handler.Filename)
	filename := strings.Join([]string{md5Filename, ext}, "")
	// Check if the file exists
	if _, err = os.Stat(filepath.Join(staticFilePosition, filename)); os.IsNotExist(err) {
		// Create file and write to file
		f, err := os.Create(filepath.Join(staticFilePosition, filename))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
		defer f.Close()
		if _, err = f.Write(content); err != nil {
			w.WriteHeader(500)
			return
		}
	}
	fileInfo := types.FileInfo{
		Filename: filename,
	}
	resBody, _ := json.Marshal(fileInfo)
	w.Write(resBody)
}

func GetUserListHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	user := model.GetUserByID(intID)
	if user.ID == -1 {
		w.WriteHeader(204)
	}
	layout := "2006-01-02 15:04"
	stringTime := user.RegisterTime.Format(layout)
	jsonRet := types.PCUserDetailedInfo{user.Name, user.Email, user.Logo, user.Evidence, user.Info, user.Account, stringTime}
	byteRet, err := json.Marshal(jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRet)
}

func VerifyPCUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	verify := r.FormValue("verify")
	intID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	intVerify, err := strconv.Atoi(verify)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if intVerify != 0 && intVerify != 1 {
		w.WriteHeader(400)
		return
	}

	// Get body message
	refuseMessage := ""
	type rejectMsg struct {
		refuseInfo string
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	jsonBody := new(rejectMsg)
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		fmt.Println(err)
	} else {
		refuseMessage = jsonBody.refuseInfo
	}

	// Update db, including time, password, account and status
	user := model.GetUserByID(intID)
	if user.ID == -1 {
		w.WriteHeader(204)
		return
	}
	password := getPassword(strconv.Itoa(user.ID), GeneratePassword(12))
	err = model.VerifyUser(intID, intVerify, user.Email, password, time.Now())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(204)
		return
	}
	// Send message to the mq
	subject := "恭喜，您的注册请求已通过"
	if intVerify == 0 {
		subject = "很抱歉，您的注册请求未被通过"
	}
	content := refuseMessage
	if intVerify == 1 {
		content = fmt.Sprintf("您的登录账户信息为: %s\r\n您的登录密码为:%s\r\n感谢您使用中大活动", user.Email, password)
	}
	msg := types.EmailContent{"admin@sysuactivity.com", user.Email, subject, content}
	byteContent, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	ok := WriteMessageQueue("email", byteContent)
	if ok {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func GetPCUserListHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stringType := r.FormValue("type")
	intType, err := strconv.Atoi(stringType)
	if err != nil || (intType != 0 && intType != 1) {
		w.WriteHeader(400)
		return
	}

	type pcuserList struct {
		content []types.PCUserListInfo
	}

	// Get user list from the db
	retContent := make([]types.PCUserListInfo, 0)
	if intType == 0 {
		// Get only the verified user
		userList := model.GetUserList(1)
		if userList == nil {
			w.WriteHeader(500)
			return
		}
		if len(userList) == 0 {
			w.WriteHeader(204)
			return
		}
		for _, v := range userList {
			layout := "2006-01-02 15:04"
			stringTime := v.RegisterTime.Format(layout)
			tmp := types.PCUserListInfo{v.ID, v.Name, v.Logo, v.Verified, stringTime}
			retContent = append(retContent, tmp)
		}
	} else {
		// Get all of the user
		userList := model.GetUserList(1)
		tmp := model.GetUserList(0)
		if userList == nil || tmp == nil {
			w.WriteHeader(500)
			return
		}
		userList = append(userList, tmp...)
		if len(userList) == 0 {
			w.WriteHeader(204)
			return
		}
		// Sort list
		sort.Sort(model.SortablePCUserList(userList))
		for _, v := range userList {
			layout := "2006-01-02 15:04"
			stringTime := v.RegisterTime.Format(layout)
			tmp := types.PCUserListInfo{v.ID, v.Name, v.Logo, v.Verified, stringTime}
			retContent = append(retContent, tmp)
		}
	}
	// Write content back to the response
	jsonRet := pcuserList{retContent}
	byteRet, err := json.Marshal(jsonRet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(byteRet)
}
