package controllers

import (
	"github.com/gorilla/mux"
	"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"context"
	"net/http"
	"io"
)

type BlockEditor struct {
	TmplEngine *templates.Engine
	StoryService *services.Story
	StoryDatabaseService *services.StoryDatabase
}

func blockEditorRetrieveState(ctx context.Context) *BlockEditor {
	return &BlockEditor{
		TmplEngine: ctx.Value(templates.EngineCtxKey).(*templates.Engine),
		StoryDatabaseService: ctx.Value(services.StoryDatabaseCtxKey).(*services.StoryDatabase),
	}
}

func getBlockEditorList(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	story := state.StoryDatabaseService.LockForRead("default")
	w.Header().Add("HX-Trigger", "evtStoryPblockListUpdate")
	state.TmplEngine.BlockEditorList(w, story)
}

func getBlockEditorForm(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])
	block := story.GetPromptBlock(vars["blockName"])
	state.TmplEngine.BlockEditorForm(w, story.Name, block, "")
}

func postBlockEditorForm(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForWrite(vars["storyName"])
	block := story.GetPromptBlock(vars["blockName"])

	input := models.PromptBlock{
		Name: r.PostFormValue("name"),
		Role: models.PromptRole(r.PostFormValue("role")),
		Text: r.PostFormValue("text"),
	}

	existing := story.GetPromptBlock(input.Name)
	if existing != nil && block != existing {
		io.WriteString(w, "A block of the same name exists already!")
		w.Header().Add("HX-Retarget", "this .info")
		w.Header().Add("HX-Reswap", "innerHTML")
		return
	}

	if block != nil {
		block.Name = input.Name
		block.Role = input.Role
		block.Text = input.Text
		state.TmplEngine.BlockEditorForm(w, story.Name, block, "Block updated!")
	} else {
		_ = story.AddPromptBlock(input)
		block = story.GetPromptBlock(input.Name)
		state.TmplEngine.BlockEditorForm(w, story.Name, block, "New block created!")
	}

}

func InstallBlockEditorController(router *mux.Router) {
	r := router.PathPrefix("/story/{storyName}/").Subrouter()
	r.HandleFunc("/blockEditor/list", getBlockEditorList).Methods("GET")
	r.HandleFunc("/blockEditor/edit/{blockName}", getBlockEditorForm).Methods("GET")
	r.HandleFunc("/blockEditor/edit/", getBlockEditorForm).Methods("GET")
	r.HandleFunc("/blockEditor/edit/{blockName}", postBlockEditorForm).Methods("POST")
	r.HandleFunc("/blockEditor/edit/", postBlockEditorForm).Methods("POST")
}

