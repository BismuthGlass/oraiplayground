package services

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"crow/orai"
	"errors"
)

// TODO: Add locks to all the calls

const StoryDatabaseCtxKey = utils.CtxKey("ServiceStoryDatabase")

type StoryDatabase struct {
	defaultModelSettings models.ModelSettings
	Stories []models.Story
}

func NewStoryDatabase() *StoryDatabase {
	defaultModelSettings := models.ModelSettings{
		Model: "lizpreciatior/lzlv-70b-fp16-hf",
		Template: "alpaca",
	}
	defaultModelSettings.SetParameters(orai.DefaultAiParameters())
	
	storyDatabase := StoryDatabase{
		defaultModelSettings: defaultModelSettings,
	}
	_ = storyDatabase.NewStory("default", "default story")
	return &storyDatabase
}

func (db *StoryDatabase) GetStory(name string) *models.Story {
	for i, _ := range db.Stories {
		if db.Stories[i].Name == name {
			return &db.Stories[i]
		}
	}
	return nil
}

func (db *StoryDatabase) NewStory(name string, description string) error {
	existingStory := db.GetStory(name)
	if existingStory != nil {
		return errors.New("exists")
	}

	db.Stories = append(db.Stories, models.NewStory(name, description, models.StoryModeInstruct, db.defaultModelSettings))
	return nil
}

func (db *StoryDatabase) LockForRead(name string) *models.Story {
	for i, _ := range db.Stories {
		if db.Stories[i].Name == name {
			// TODO: Lock story for read
			return &db.Stories[i]
		}
	}
	return nil
}

func (db *StoryDatabase) LockForWrite(name string) *models.Story {
	for i, _ := range db.Stories {
		if db.Stories[i].Name == name {
			// TODO: Lock story for write
			return &db.Stories[i]
		}
	}
	return nil
}

