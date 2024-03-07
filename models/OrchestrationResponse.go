package models

type OrchestrationResponse struct {
	Response []Result `json:"response"`
}

type Result struct {
	Provider   SystemDefinition  `json:"provider"`
	ServiceUri string            `json:"serviceUri"`
	Metadata   map[string]string `json:"metadata"`
}
