package router

import (
	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/controller"
	"github.com/sysu-activitypluspc/service-end/middleware"
	"github.com/urfave/negroni"
)

func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	s := negroni.Classic()
	s.Use(middleware.ValidUserMiddleWare{})
	

	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", controller.AddActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", controller.ModifyActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", controller.DeleteActivityHandler).Methods("DELETE")
	act.HandleFunc("", controller.VerifyActivityHandler).Methods("PUT")
	act.HandleFunc("", controller.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/{id}", controller.ShowActivityDetailHandler).Methods("GET")
	act.HandleFunc("/{clubId}/list", controller.ShowActivitiesListByClubHandler).Methods("GET")
	act.HandleFunc("/{clubId}/status", controller.GetNumberOfActStatusByClub).Methods("GET")
	act.HandleFunc("/{actid}", controller.CloseActivityApply).Methods("PUT")

	session := r.PathPrefix("/session").Subrouter()
	session.HandleFunc("", controller.LoginHandler).Methods("POST")

	pcuser := r.PathPrefix("/pcusers").Subrouter()
	pcuser.HandleFunc("", controller.SignUpHandler).Methods("POST")
	pcuser.HandleFunc("/{id}", controller.GetPCUserDetailHandler).Methods("GET")
	pcuser.HandleFunc("", controller.VerifyPCUserHandler).Methods("PUT")
	pcuser.HandleFunc("", controller.GetPCUserListHandler).Methods("GET")

	message := r.PathPrefix("/message").Subrouter()
	message.HandleFunc("", controller.AddMessageHandler).Methods("POST")
	message.HandleFunc("", controller.ShowMessagesListHandler).Methods("GET")
	message.HandleFunc("/{id}", controller.ShowMessageDetailHandler).Methods("GET")
	message.HandleFunc("/{id}", controller.DeleteMessageHandler).Methods("DELETE")
	
	r.HandleFunc("/images", controller.UploadImageHandler).Methods("POST")

	actApply := r.PathPrefix("/actApply").Subrouter()
	actApply.HandleFunc("", controller.ListActivityApplyHandler).Methods("GET")
	actApply.HandleFunc("", controller.DeleteActivityApplyHandler).Methods("DELETE")
	

	s.UseHandler(r)
	return s
}
