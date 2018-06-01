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

	session := r.PathPrefix("/session").Subrouter()
	session.HandleFunc("", controller.LoginHandler).Methods("POST")

	pcuser := r.PathPrefix("/pcusers").Subrouter()
	pcuser.HandleFunc("", controller.SignUpHandler).Methods("POST")
	pcuser.HandleFunc("/{id}", controller.GetPCUserDetailHandler).Methods("GET")
	pcuser.HandleFunc("", controller.VerifyPCUserHandler).Methods("PUT")
	pcuser.HandleFunc("", controller.GetPCUserListHandler).Methods("GET")

	r.HandleFunc("/images", controller.UploadImageHandler).Methods("POST")

	s.UseHandler(r)
	return s
}
