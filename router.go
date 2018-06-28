package main

import (
	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/middleware"
	"github.com/sysu-activitypluspc/service-end/service"
	"github.com/urfave/negroni"
)

/*
* /images 0
* /pcusers post 0
* /session 0
* /act post 0
* /act get put 2
* /act/ delete 2
* /message 2
* /pcusers post 0
* /pcusers put get 2
* /pcusers/ 2
* /act/.../list clubid == userid
* .../status clubid == userid
* /act/ post put act is published by user
* /actApply act is published by user
*/

func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	s := negroni.Classic()
	s.Use(middleware.ParseToken)

	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", service.AddActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", service.ModifyActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", service.DeleteActivityHandler).Methods("DELETE")
	act.HandleFunc("", service.VerifyActivityHandler).Methods("PUT")
	act.HandleFunc("", service.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/{id}", service.ShowActivityDetailHandler).Methods("GET")
	act.HandleFunc("/{clubId}/list", service.ShowActivitiesListByClubHandler).Methods("GET")
	act.HandleFunc("/{clubId}/status", service.GetNumberOfActStatusByClubHandler).Methods("GET")
	act.HandleFunc("/{actid}", service.CloseActivityHandler).Methods("PUT")

	session := r.PathPrefix("/session").Subrouter()
	session.HandleFunc("", service.LoginHandler).Methods("POST")

	pcuser := r.PathPrefix("/pcusers").Subrouter()
	pcuser.HandleFunc("", service.SignUpHandler).Methods("POST")
	pcuser.HandleFunc("/{id}", service.GetPCUserDetailHandler).Methods("GET")
	pcuser.HandleFunc("", service.VerifyPCUserHandler).Methods("PUT")
	pcuser.HandleFunc("", service.GetPCUserListHandler).Methods("GET")

	message := r.PathPrefix("/message").Subrouter()
	message.HandleFunc("", service.AddMessageHandler).Methods("POST")
	message.HandleFunc("", service.ShowMessagesListHandler).Methods("GET")
	message.HandleFunc("/{id}", service.ShowMessageDetailHandler).Methods("GET")
	message.HandleFunc("/{id}", service.DeleteMessageHandler).Methods("DELETE")

	r.HandleFunc("/images", service.UploadImageHandler).Methods("POST")

	actApply := r.PathPrefix("/actApply").Subrouter()
	actApply.HandleFunc("", service.ListActivityApplyHandler).Methods("GET")
	actApply.HandleFunc("", service.DeleteActivityApplyHandler).Methods("DELETE")

	s.UseHandler(r)
	return s
}
