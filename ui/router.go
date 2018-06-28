package ui

import (
	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/service"
	"github.com/sysu-activitypluspc/service-end/middleware"
	"github.com/urfave/negroni"
)

func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	s := negroni.Classic()
	s.Use(middleware.ValidUserMiddleWare{})
	

	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", service.AddActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", service.ModifyActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", service.DeleteActivityHandler).Methods("DELETE")
	act.HandleFunc("", service.VerifyActivityHandler).Methods("PUT")
	act.HandleFunc("", service.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/{id}", service.ShowActivityDetailHandler).Methods("GET")
	act.HandleFunc("/{clubId}/list", service.ShowActivitiesListByClubHandler).Methods("GET")
	act.HandleFunc("/{clubId}/status", service.GetNumberOfActStatusByClub).Methods("GET")
	act.HandleFunc("/{actid}", service.CloseActivityApply).Methods("PUT")

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
