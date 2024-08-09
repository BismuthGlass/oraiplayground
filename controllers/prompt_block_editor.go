package controllers

import (
	"github.com/gorilla/mux"
	//"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"context"
	"net/http"
)

type PromptBlockEditor struct {
	TmplEngine *templates.Engine
	StoryService *services.Story
	StoryDatabaseService *services.StoryDatabase
}

func promptBlockEditorRetrieveState(ctx context.Context) *PromptBlockEditor {
	return &PromptBlockEditor{
		TmplEngine: ctx.Value(templates.EngineCtxKey).(*templates.Engine),
		StoryDatabaseService: ctx.Value(services.StoryDatabaseCtxKey).(*services.StoryDatabase),
	}
}

func getPromptBlockEditorList(w http.ResponseWriter, r *http.Request) {
	state := promptBlockEditorRetrieveState(r.Context())
	story := state.StoryDatabaseService.LockForRead("default")
	w.Header().Add("HX-Trigger", "evtStoryPblockListUpdate")
	state.TmplEngine.PromptBlockEditorTable(w, story)
}

func InstallPromptBlockEditorController(router *mux.Router) {
	r := router.PathPrefix("/story/{storyName}/").Subrouter()
	r.HandleFunc("/promptBlockEditor/list", getPromptBlockEditorList).Methods("GET")
}

