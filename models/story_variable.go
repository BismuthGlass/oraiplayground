package models

type StoryVariableType string
const (
	StoryVariableTypeString = "string"
	StoryVariableTypeNumber = "number"
	StoryVariableTypeCode = "code"
)

type StoryVariable struct {
	Name string
	Type StoryVariableType
	StringValue string
	NumberValue int
}

func NewNumberStoryVariable(name string, val int) StoryVariable {
	return StoryVariable{
		Name: name,
		Type: StoryVariableTypeNumber,
		NumberValue: val,
	}
}

func NewStringStoryVariable(name string, val string) StoryVariable {
	return StoryVariable{
		Name: name,
		Type: StoryVariableTypeString,
		StringValue: val,
	}
}

func NewCodeStoryVariable(name string, val string) StoryVariable {
	return StoryVariable{
		Name: name,
		Type: StoryVariableTypeCode,
		StringValue: val,
	}
}

