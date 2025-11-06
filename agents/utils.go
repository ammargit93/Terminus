package agents

import "encoding/json"

type Argument struct {
	ActionName string   `json:"actionName"`
	ArgNames   []string `json:"argNames"`
}
type ResponseMessage struct {
	Action  []string   `json:"action"`
	Message string     `json:"message"`
	Code    string     `json:"code"`
	Args    []Argument `json:"args"`
	ToolUse bool       `json:"tooluse"`
}

func ParseJSON(content string) ResponseMessage {
	var v ResponseMessage
	json.Unmarshal([]byte(content), &v)
	return v
}
