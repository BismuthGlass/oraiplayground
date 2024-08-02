package controllers

import (
	"context"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PromptBlock struct {
	TmplEngine *templates.Engine
	StoryService *services.Story	
}

func promptBlockRetrieveState(ctx context.Context) PromptBlock {
	return PromptBlock{
		TmplEngine: ctx.Value(templates.EngineCtxKey).(*templates.Engine),
		StoryService: ctx.Value(services.StoryCtxKey).(*services.Story),
	}
}

func putPromptBlockList(w http.ResponseWriter, r *http.Request) {
	state := promptBlockRetrieveState(r.Context())
	r.ParseForm()

	moveTargetBlock := r.PostFormValue("moveTargetBlock")
	moveDestinationGroup := r.PostFormValue("moveDestinationGroup")
	moveDestinationBlock := r.PostFormValue("moveDestinationBlock")

	err := state.StoryService.Story.PromptRoot.ChangeLocation(
		moveTargetBlock,
		moveDestinationGroup,
		moveDestinationBlock,
	)
	if err != nil {
		log.Println(err)
	}

	ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	state.TmplEngine.PromptBlockList(w, &ctx)
}

func InstallPromptBlockController(router *mux.Router) {
	router.HandleFunc("/promptBlockList", putPromptBlockList).Methods("PUT")
}