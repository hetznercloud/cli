package schema

import "time"

// Action defines the schema of an action.
type Action struct {
	ID       int          `json:"id"`
	Status   string       `json:"status"`
	Command  string       `json:"command"`
	Progress int          `json:"progress"`
	Started  time.Time    `json:"started"`
	Finished time.Time    `json:"finished"`
	Error    *ActionError `json:"error"`
}

// ActionError defines the schema of an error embedded
// in an action.
type ActionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ActionGetResponse is the schema of the response when
// retrieving a single action.
type ActionGetResponse struct {
	Action Action `json:"action"`
}
