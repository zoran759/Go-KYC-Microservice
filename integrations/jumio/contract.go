package jumio

// Jumio performNetverify API endpoints.
const (
	USbaseURL = "https://netverify.com/api/netverify/v2"
	EUbaseURL = "https://lon.netverify.com/api/netverify/v2"
)

// Config holds configuration settings for the service.
type Config struct {
	BaseURL string
	Token   string
	Secret  string
}
