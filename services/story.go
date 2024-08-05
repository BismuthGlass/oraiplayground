package services

import (
	"crow/oraiplayground/utils"
)

const StoryCtxKey = utils.CtxKey("ServiceStory")

type Story struct {
	AiServer *AiServer
}

func NewStory(aiServer *AiServer) *Story {
	return &Story{
		AiServer: aiServer,
	}
}

func (s *Story) GetPrompt() string {
	//return s.Story.PromptRoot.IntoPrompt(s.Story.Settings.Template)
	return ""
}

func (s *Story) RequestPrompt(continueCue string) AiServiceRequestClientInfo {
	//return s.AiServer.IssueRequest(
	//	s.Story.Settings.Model, 
	//	s.Story.Settings.GetParameters(),
	//	BlocksIntoPrompt(story.Settings.Template) + continueCue,
	//)
	return AiServiceRequestClientInfo{}
}

func (s *Story) GetPromptRequest(id int64) AiServiceResponse {
	return <-s.AiServer.RequestChannel(id)
}

func (s *Story) CancelPromptRequest(id int64) error {
	return s.AiServer.CancelRequest(id)
}

//func (s *Story) AddMessage(role orai.ChatRole, message string) {
//	s.Story.ChatMessages = append(s.Story.ChatMessages, orai.ChatMessage{
//		Role: role,
//		Message: message,
//	})
//}

//func (s *Story) RequestMessage(role orai.ChatRole, message string) AiServiceRequestClientInfo {
//	blockMessages := s.Story.PromptRoot.IntoMessages()
//	var cue string
//	switch role {
//	case orai.ChatRoleUser:
//		cue = s.Story.PromptRoot.GatherRole(models.PromptRoleUserCue) + message
//	case orai.ChatRoleAssistant:
//		cue = s.Story.PromptRoot.GatherRole(models.PromptRoleAssistantCue) + message
//	default:
//		cue = message
//	}
//	prompt := orai.GenerateLlama3Prompt(blockMessages, role, cue)
//	return s.AiServer.IssueRequest(
//		s.Story.Settings.Model, 
//		s.Story.Settings.GetParameters(),
//		prompt,
//	)
//}
