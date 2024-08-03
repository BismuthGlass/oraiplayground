package controllers

import (
	"context"
	"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"crow/oraiplayground/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Story struct {
	TmplEngine *templates.Engine
	StoryService *services.Story
}

func storyRetrieveState(ctx context.Context) *Story {
	return &Story{
		TmplEngine: ctx.Value(templates.EngineCtxKey).(*templates.Engine),
		StoryService: ctx.Value(services.StoryCtxKey).(*services.Story),
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())

	ctx := templates.Story{
		Mode: state.StoryService.Story.Settings.Mode,
		Settings: state.StoryService.SettingsTmplModel(nil),
		PromptBlockList: state.StoryService.PromptListComponent.IntoTmplModel(),
	}

	_ = state.TmplEngine.StoryPage(w, &ctx)
}

func postSettings(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()
	
	err := state.StoryService.Story.Settings.ParseFormData(r.PostForm)
	ctx := state.StoryService.SettingsTmplModel(err)
	state.TmplEngine.Settings(w, &ctx)
}

func getBlockEditor(w http.ResponseWriter, r *http.Request) {
	// TODO: What about locks?
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	block, exists := state.StoryService.Story.PromptRoot.GetBlock(r.Form.Get("block"))

	blockEditor := templates.PromptBlockEditor{
		Mode: state.StoryService.Story.Settings.Mode,
		EditorId: r.Form.Get("eid"),
		SectionOptions: models.PromptBlockSectionOptions(),
		RoleOptions: models.PromptBlockRoleOptions(),
	}
	if exists {
		blockEditor.Name = block.Name
		blockEditor.IsGroup = block.IsGroup
		blockEditor.Text = block.Text
		blockEditor.Prefix = block.Prefix
		blockEditor.Suffix = block.Suffix
		utils.SetSelection(blockEditor.RoleOptions, string(block.Role))
		utils.SetSelection(blockEditor.SectionOptions, string(block.Section))
	}

	_ = state.TmplEngine.PromptBlockEditor(w, &blockEditor)
}

func putBlockEditor(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	updateBlock := models.PromptBlock{
		Name: r.Form.Get("block"),
		Section: models.PromptSection(r.PostFormValue("section")),
		Text: r.PostFormValue("text"),
		Prefix: r.PostFormValue("prefix"),
		Suffix: r.PostFormValue("suffix"),
	}

	block := state.StoryService.Story.PromptRoot.UpdateBlock(updateBlock)
	if block == nil {
		return
	}
}

func postPromptBlock(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	newBlock := models.PromptBlock{
		Name: r.PostFormValue("name"),
		IsGroup: r.PostFormValue("type") == "group",
	}

	// TODO: Handle this error
	_ = state.StoryService.NewPromptBlock(newBlock)

	ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	state.TmplEngine.PromptBlockList(w, &ctx)
}

func deletePromptBlock(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	blockNames := r.PostFormValue("blocks")
	blockNameList := strings.Split(blockNames, ",")
	state.StoryService.DeletePromptBlocks(blockNameList)

	ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	state.TmplEngine.PromptBlockList(w, &ctx)
}

func postAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	err := r.ParseMultipartForm(10000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	continueValue := r.PostFormValue("continue")

	reqInfo := state.StoryService.RequestPrompt(continueValue)
	jsonData, err := json.Marshal(&reqInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func getAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := state.StoryService.GetPromptRequest(id)
	jsonData, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteAiRequest(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = state.StoryService.CancelPromptRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type PromptResponse struct {
	Prompt string `json:"prompt"`
}

func getPromptPreview(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	prompt := state.StoryService.GetPrompt()

	res := PromptResponse{
		Prompt: prompt,
	}
	jsonData, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func putActivatePromptBlock(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	name := r.PostFormValue("name")
	active := r.PostFormValue("active") == "on"
	state.StoryService.ActivatePromptBlock(name, active)

	ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	state.TmplEngine.PromptBlockList(w, &ctx)
}

func InstallStoryController(router *mux.Router) {
	router.HandleFunc("/", getIndex).Methods("GET")
	router.HandleFunc("/settings", postSettings).Methods("POST")
	router.HandleFunc("/blockEditor", getBlockEditor).Methods("GET")
	router.HandleFunc("/blockEditor", putBlockEditor).Methods("PUT")
	router.HandleFunc("/promptBlock", postPromptBlock).Methods("POST")
	router.HandleFunc("/promptBlock/delete", deletePromptBlock).Methods("POST")
	router.HandleFunc("/aiRequest", postAiRequest).Methods("POST")
	router.HandleFunc("/aiRequest", getAiRequest).Methods("GET")
	router.HandleFunc("/aiRequest", deleteAiRequest).Methods("DELETE")
	router.HandleFunc("/promptPreview", getPromptPreview).Methods("GET")
	router.HandleFunc("/activatePromptBlock", putActivatePromptBlock).Methods("PUT")
}
