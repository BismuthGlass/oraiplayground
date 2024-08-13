package controllers

import (
	"github.com/gorilla/mux"
	"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"context"
	"net/http"
	"io"
	"log"
)

type BlockEditor struct {
	StoryService *services.Story
	StoryDatabaseService *services.StoryDatabase
}

func blockEditorRetrieveState(ctx context.Context) *BlockEditor {
	return &BlockEditor{
		StoryDatabaseService: ctx.Value(services.StoryDatabaseCtxKey).(*services.StoryDatabase),
	}
}

func getBlockEditorList(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	story := state.StoryDatabaseService.LockForRead("default")
	w.Header().Add("HX-Trigger", "evtStoryPblockListUpdate")
	templates.BlockEditorList(w, story)
}

func getBlockEditorForm(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])
	block := story.GetPromptBlock(vars["blockName"])
	templates.BlockEditorForm(w, story.Name, block, "")
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
		w.Header().Add("HX-Retarget", "previous .info")
		w.Header().Add("HX-Reswap", "innerHTML")
		io.WriteString(w, "A block of the same name exists already!")
		return
	}

	w.Header().Add("HX-Trigger", "updateEditorBlockList")

	if block != nil {
		block.Name = input.Name
		block.Role = input.Role
		block.Text = input.Text
		templates.BlockEditorForm(w, story.Name, block, "Block updated!")
	} else {
		_ = story.AddPromptBlock(input)
		block = story.GetPromptBlock(input.Name)
		templates.BlockEditorForm(w, story.Name, block, "New block created!")
	}

}

func putBlockEditorFavorite(w http.ResponseWriter, r *http.Request) {
	state := blockEditorRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForWrite(vars["storyName"])
	err := story.TogglePromptBlockFavorite(vars["blockName"])
	if err != nil {
		log.Fatalf(err.Error())
	}
	templates.BlockEditorList(w, story)
}

func InstallBlockEditorController(router *mux.Router) {
	r := router.PathPrefix("/story/{storyName}/blockEditor/").Subrouter()
	r.HandleFunc("/list", getBlockEditorList).Methods("GET")
	r.HandleFunc("/edit/{blockName}", getBlockEditorForm).Methods("GET")
	r.HandleFunc("/edit/", getBlockEditorForm).Methods("GET")
	r.HandleFunc("/edit/{blockName}", postBlockEditorForm).Methods("POST")
	r.HandleFunc("/edit/", postBlockEditorForm).Methods("POST")
	r.HandleFunc("/favorite/{blockName}", putBlockEditorFavorite).Methods("PUT")
}

