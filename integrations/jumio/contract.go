package jumio

// Jumio performNetverify API endpoints.
const (
	USendpoint = "https://netverify.com/api/netverify/v2/performNetverify"
	EUendpoint = "https://lon.netverify.com/api/netverify/v2/performNetverify"
)

// Config holds configuration settings for the service.
type Config struct {
	Host   string
	Token  string
	Secret string
}
