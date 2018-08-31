package identitymind

import "gitlab.com/lambospeed/kyc/integrations/identitymind/consumer"

// IdentityMind API urls for the convenience.
const (
	SandboxBaseURL    = "https://sandbox.identitymind.com/im"
	StagingBaseURL    = "https://staging.identitymind.com/im"
	ProductionBaseURL = "https://edna.identitymind.com/im"
)

// Config holds configuration settings for the service.
type Config struct {
	Host     string
	Username string
	Password string
}

var _ *Config = (*Config)(&consumer.Config{})
