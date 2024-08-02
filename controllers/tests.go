package controllers

import (
	"crow/oraiplayground/templates"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getTestHelloOOB(w http.ResponseWriter, r *http.Request) {
	templates := r.Context().Value(templates.EngineCtxKey).(*templates.Engine)

	err := templates.Template.ExecuteTemplate(w, "tests/hello-oob.html", nil)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}

func getTestList(w http.ResponseWriter, r *http.Request) {
	templates := r.Context().Value(templates.EngineCtxKey).(*templates.Engine)

	err := templates.Template.ExecuteTemplate(w, "tests/list.html", nil)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}

func InstallTestsController(router *mux.Router) {
	router.HandleFunc("/tests/hello-oob", getTestHelloOOB).Methods("GET")
	router.HandleFunc("/tests/list", getTestList).Methods("GET")
}
