package events

import "automation-hub-backend/internal/models"

type AutomationEventType string

const (
	CreateEvent AutomationEventType = "create"
	UpdateEvent AutomationEventType = "update"
	DeleteEvent AutomationEventType = "delete"
)

type AutomationEvent struct {
	Type       AutomationEventType `json:"type"`
	Automation *models.Automation  `json:"automation"`
}
