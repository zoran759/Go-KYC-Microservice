package verification

// Client represents the client of the ShuftiPro API.
// It shouldn't initialized directly, use New() constructor instead.
type Client struct {
	config Config
}

// New constructs new Client object.
func New(config Config) Client {
	return Client{
		config: config,
	}
}
