package router

import (
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
	"github.com/sysu-activitypluspc/service-end/controller"
)

func GetServer() *negroni.Negroni{
	r := mux.NewRouter()

	s := negroni.Classic()

	// TODO: Add handler here
	r.HandleFunc("/act", controller.AddActivityHandler).Methods("POST")

	s.UseHandler(r)
	return s
}