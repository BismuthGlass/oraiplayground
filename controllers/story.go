package controllers

import (
	"context"
	"crow/oraiplayground/models"
	"crow/oraiplayground/services"
	"crow/oraiplayground/templates"
	"net/http"

	"github.com/gorilla/mux"
)

type Story struct {
	StoryDatabaseService *services.StoryDatabase
}

func storyRetrieveState(ctx context.Context) *Story {
	return &Story{
		StoryDatabaseService: ctx.Value(services.StoryDatabaseCtxKey).(*services.StoryDatabase),
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

func getBlockEditor(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//r.ParseForm()

	//story := state.StoryDatabaseService.LockForRead("default")
	//block := story.GetPromptBlock(r.Form.Get("block"))

	//blockEditor := templates.PromptBlockEditor{
	//	Mode: story.Mode,
	//	EditorId: r.Form.Get("eid"),
	//	SectionOptions: models.PromptBlockSectionOptions(),
	//	RoleOptions: models.PromptBlockRoleOptions(),
	//}
	//if block != nil {
	//	blockEditor.Name = block.Name
	//	blockEditor.Text = block.Text
	//	utils.SetSelection(blockEditor.RoleOptions, string(block.Role))
	//}

	//_ = templates.PromptBlockEditor(w, &blockEditor)
}

func putBlockEditor(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	r.ParseForm()

	story := state.StoryDatabaseService.LockForWrite("default")

	updateBlock := models.PromptBlock{
		Name: r.Form.Get("block"),
		Role: models.PromptRole(r.PostFormValue("role")),
		Text: r.PostFormValue("text"),
	}

	// TODO
	err := story.UpdatePromptBlock(r.Form.Get("block"), updateBlock)
	if err != nil {
		return
	}
}

func postPromptBlock(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//r.ParseForm()

	//newBlock := models.PromptBlock{
	//	Name: r.PostFormValue("name"),
	//	IsGroup: r.PostFormValue("type") == "group",
	//}

	//_ = state.StoryService.NewPromptBlock(newBlock)

	//ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	//templates.PromptBlockList(w, &ctx)
}

func deletePromptBlock(w http.ResponseWriter, r *http.Request) {
//	state := storyRetrieveState(r.Context())
//	r.ParseForm()
//
//	blockNames := r.PostFormValue("blocks")
//	blockNameList := strings.Split(blockNames, ",")
//	state.StoryService.DeletePromptBlocks(blockNameList)
//
//	ctx := state.StoryService.PromptListComponent.IntoTmplModel()
//	templates.PromptBlockList(w, &ctx)
}

func postAiRequest(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//err := r.ParseMultipartForm(10000)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//continueValue := r.PostFormValue("continue")

	//reqInfo := state.StoryService.RequestPrompt(continueValue)
	//jsonData, err := json.Marshal(&reqInfo)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//r.Header.Add("Content-Type", "application/json")
	//w.Write(jsonData)
}

func getAiRequest(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//idStr := r.URL.Query().Get("id")
	//id, err := strconv.ParseInt(idStr, 10, 64)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	//response := state.StoryService.GetPromptRequest(id)
	//jsonData, err := json.Marshal(&response)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//r.Header.Add("Content-Type", "application/json")
	//w.Write(jsonData)
}

func deleteAiRequest(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//idStr := r.URL.Query().Get("id")
	//id, err := strconv.ParseInt(idStr, 10, 64)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	//err = state.StoryService.CancelPromptRequest(id)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//}
}

type PromptResponse struct {
	Prompt string `json:"prompt"`
}

func getPromptPreview(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//prompt := state.StoryService.GetPrompt()

	//res := PromptResponse{
	//	Prompt: prompt,
	//}
	//jsonData, err := json.Marshal(&res)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//r.Header.Add("Content-Type", "application/json")
	//w.Write(jsonData)
}

func putActivatePromptBlock(w http.ResponseWriter, r *http.Request) {
	//state := storyRetrieveState(r.Context())
	//r.ParseForm()

	//name := r.PostFormValue("name")
	//active := r.PostFormValue("active") == "on"
	//state.StoryService.ActivatePromptBlock(name, active)

	//ctx := state.StoryService.PromptListComponent.IntoTmplModel()
	//templates.PromptBlockList(w, &ctx)
}

func getPromptInfo(w http.ResponseWriter, r *http.Request) {
	state := storyRetrieveState(r.Context())
	vars := mux.Vars(r)
	story := state.StoryDatabaseService.LockForRead(vars["storyName"])
	templates.PromptInfo(w, story)
}

func InstallStoryController(router *mux.Router) {
	router.HandleFunc("/", getIndex).Methods("GET")
	r := router.PathPrefix("/story/{storyName}/").Subrouter()
	r.HandleFunc("/promptInfo", getPromptInfo).Methods("GET")
}
