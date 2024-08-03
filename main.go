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
	tmplEngine := &templates.Engine{
		Template: tmpl,
	}
	
	aiServerService := services.NewAiServer()

	storyService := services.NewStory(&aiServerService)

	go aiServerService.Run()

	// appState.AIService = NewAIService()

	rt := mux.NewRouter()
	rt.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("data/static"))))
	controllers.InstallStoryController(rt)
	controllers.InstallPromptBlockController(rt)
	controllers.InstallTestsController(rt)

	// Start services
	// go appState.AIService.Run()

	httpConfig := http.Server{
		Addr: ":8080",
		Handler: rt,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.Background()
			ctx = context.WithValue(ctx, templates.EngineCtxKey, tmplEngine)
			ctx = context.WithValue(ctx, services.StoryCtxKey, storyService)
			return ctx
		},
	}

	err = httpConfig.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
