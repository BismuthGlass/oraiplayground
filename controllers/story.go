package controllers

import (
	"context"
	"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"net/http"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/mux"
)

type Story struct {
	StoryDatabaseService *services.StoryDatabase
	AiServerService *services.AiServer
}

func storyRetrieveState(ctx context.Context) *Story {
	return &Story{
		StoryDatabaseService: ctx.Value(services.StoryDatabaseCtxKey).(*services.StoryDatabase),
		AiServerService: ctx.Value(services.AiServerCtxKey).(*services.AiServer),
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	story := state.StoryDatabaseService.LockForRead("default")
	_ = templates.StoryPage(w, story)
}

func postSettings(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//r.ParseForm()
	//
	//err := state.StoryService.Story.Settings.ParseFormData(r.PostForm)
	//ctx := templates.NewStorySettings(&state.StoryService.Story, err)
	//templates.Settings(w, &ctx)
}

func getPlaygroundList(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])
	templates.PlaygroundBlockList(w, story)
}

func getPromptInfo(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])
	templates.PromptInfo(w, story)
}

func postAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])

	prompt := models.BlocksIntoPrompt(story.GetEnabledBlocks(), story.ModelSettings.Template)

	reqInfo := state.AiServerService.IssueRequest(
		story.ModelSettings.Model,
		story.ModelSettings.GetParameters(),
		prompt,
	)

	jsonData, err := json.Marshal(&reqInfo)
	if err != nil {
		log.Fatalf(err.Error())
	}

	r.Header.Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func getAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["requestId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch := state.AiServerService.RequestChannel(id)
	data := <-ch
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	r.Header.Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["requestId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = state.AiServerService.CancelRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func InstallStoryController(router *mux.Router) {
	router.HandleFunc("/", getIndex).Methods("GET")
	r := router.PathPrefix("/story/{storyName}/").Subrouter()
	r.HandleFunc("/promptInfo", getPromptInfo).Methods("GET")
	r.HandleFunc("/playground/list", getPlaygroundList).Methods("GET")
	r.HandleFunc("/gen", postAiRequest).Methods("POST")
	r.HandleFunc("/gen/{requestId}", getAiRequest).Methods("GET")
	r.HandleFunc("/gen/{requestId}", deleteAiRequest).Methods("DELETE")
}

