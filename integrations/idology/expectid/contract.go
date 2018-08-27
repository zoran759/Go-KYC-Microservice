package expectid

const (
	// APIendpoint holds IDology ExpectID® API Endpoint.
	APIendpoint = "https://web.idologylive.com/api/idiq.svc"
)

// Config holds configuration settings for the IDology ExpectID® API client.
type Config struct {
	Host             string
	Username         string
	Password         string
	UseSummaryResult bool
}
