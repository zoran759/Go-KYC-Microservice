package model

// NewParticipant represents a new participant without data.
// A new participant creation has to be requested first then send details if succeed.
type NewParticipant struct {
	Email string `json:"email"`
	// TODO: perhaps, need to implement "addresses" field strings preparation (!!!check the apiary for the format!!!).
	Addresses []string `json:"addresses,omitempty"`
	UUID      string   `json:"uuid,omitempty"`
}

// NewParticipantResponse represents the response on the request of adding of the new participant.
type NewParticipantResponse struct {
	UUID string `json:"uuid"`
}

// ParticipantDetails represents an individual participant data.
// They have to be send after successful creation of the new participant.
type ParticipantDetails struct {
	// TODO: write this.
}
