package router

import (
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

func GetServer() *negroni.Negroni{
	r := mux.NewRouter()

	s := negroni.Classic()

	// TODO: Add handler here

	s.UseHandler(r)
	return s
}