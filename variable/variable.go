package variable

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/telemac/goutils/natsevents"
	"time"
)

// Variable represents represents a variable value
type Variable struct {
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"ts,omitempty"`
	Comment   string      `json:"comment,omitempty"`
}

// Variables is an array of variable values
type Variables []Variable

// cloud event variable set type
const CESetType = "com.plugis.variable.set"

func NewVariableSetEvent(variables Variables) *cloudevents.Event {
	return natsevents.NewEvent("", CESetType, variables)
}
