package router

import (
	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/controller"
	"github.com/urfave/negroni"
)

func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	s := negroni.Classic()

	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", controller.AddActivityHandler).Methods("POST")
	act.HandleFunc("/", controller.AddActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", controller.ModifyActivityHandler).Methods("POST")
	act.HandleFunc("/{actId}", controller.DeleteActivityHandler).Methods("DELETE")
	act.HandleFunc("/", controller.VerifyActivityHandler).Methods("PUT")
	act.HandleFunc("", controller.VerifyActivityHandler).Methods("PUT")

	s.UseHandler(r)
	return s
}
