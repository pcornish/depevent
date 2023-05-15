package eventmodel

type Parameter struct {
	ParameterKey   string `json:"ParameterKey"`
	ParameterValue string `json:"ParameterValue"`
}

type Message struct {
	// contains the 'Environment' field
	StackTemplateParameters []Parameter `json:"stackTemplateParameters"`

	// for example: 'SUCCEEDED'
	DeploymentStatus string `json:"deploymentStatus"`

	GitCommitSha  string `json:"gitCommitSha"`
	GitRepository string `json:"gitRepository"`
}

type RequestPayload struct {
	Time string `json:"time"`
}

type ResponsePayload struct {
	Message *Message `json:"message"`
}

type DeploymentEvent struct {
	RequestPayload  *RequestPayload  `json:"requestPayload"`
	ResponsePayload *ResponsePayload `json:"responsePayload"`
}
