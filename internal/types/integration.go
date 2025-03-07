package types

import "time"

// Transcription represents the transcription details
type Transcription struct {
	AI    string    `json:"ai"`
	Human string    `json:"human"`
	Time  time.Time `json:"time"`
}

// Call represents a single call record
type Call struct {
	CallID             string          `json:"call_id"`
	PatientPhoneNumebr string          `json:"Patient_phone_Numebr"`
	CallerPhoneNumber  string          `json:"caller_phone_number"`
	Timestamp          string          `json:"timestamp"`
	Duration           int             `json:"duration"`
	ExternalCallID     string          `json:"external_call_id"`
	QuickDescription   string          `json:"summary"`
	TenantID           string          `json:"tenant_id"`
	Transcription      []Transcription `json:"transcription"`
	TransferredTo      string          `json:"transferred_to"`
}

// CallsResponse represents the full response
type CallsResponse struct {
	Calls            []Call `json:"calls"`
	LastEvaluatedKey string `json:"last_evaluated_key"`
}
