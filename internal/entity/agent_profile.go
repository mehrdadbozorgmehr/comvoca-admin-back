package entity

import "github.com/google/uuid"

type LanguageType int

const (
	English LanguageType = iota
	French
	Spanish
)

type VoiceType int

const (
	Professional VoiceType = iota
	Friendly
	Casual
)

type AgentProfile struct {
	Id            uuid.UUID
	Language      LanguageType
	InitialScript string
	ClosingScript string
	VoiceType     string
}
