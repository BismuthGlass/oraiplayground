package main

import (
	"crow/oraiplayground/config"
	"crow/oraiplayground/controllers"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"crow/webutil"

	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()

	tmpl, err := webutil.LoadTemplates("data/templates", []string{".html"})
	if err != nil {
		fmt.Println(err)
		return
	}
	templates.Init(tmpl)

	aiServerService := services.NewAiServer()
	storyDatabaseService := services.NewStoryDatabase()
	storyService := services.NewStory(&aiServerService)

	go aiServerService.Run()

	rt := mux.NewRouter()
	rt.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("data/static"))))
	controllers.InstallStoryController(rt)
	controllers.InstallTestsController(rt)
	controllers.InstallBlockEditorController(rt)

	httpConfig := http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: rt,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.Background()
			ctx = context.WithValue(ctx, services.StoryDatabaseCtxKey, storyDatabaseService)
			ctx = context.WithValue(ctx, services.StoryCtxKey, storyService)
			return ctx
		},
	}

	fmt.Printf("Listening on port %s...\n", config.ServerPort)
	err = httpConfig.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
