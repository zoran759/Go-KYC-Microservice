package jumio

// Response defines the model for the performNetverify API response.
type Response struct {
	// Timestamp (UTC) of the response format: YYYY-MM-DDThh:mm:ss.SSSZ.
	Timestamp string `json:"timestamp"`
	// Jumio's reference number for each scan. Max. length 36.
	JumioIDScanReference string `json:"jumioIdScanReference"`
}
