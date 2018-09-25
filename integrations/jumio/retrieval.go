package jumio

import (
	"time"
)

// The endpoints to use for the scan status retrieval.
const (
	usScanStatusEndpoint = "https://netverify.com/api/netverify/v2/scans/"
	euScanStatusEndpoint = "https://lon.netverify.com/api/netverify/v2/scans/"
)

// Request timings recommendation according to https://github.com/Jumio/implementation-guides/blob/master/netverify/netverify-retrieval-api.md#usage.
var timings = [10]time.Duration{
	40 * time.Second,
	60 * time.Second,
	100 * time.Second,
	160 * time.Second,
	240 * time.Second,
	340 * time.Second,
	460 * time.Second,
	600 * time.Second,
	760 * time.Second,
	940 * time.Second,
}

// StatusResponse defines the model for the response on retrieving scan status.
type StatusResponse struct {
	// Timestamp of the response in the format YYYY-MM-DDThh:mm:ss.SSSZ.
	Timestamp string `json:"timestamp"`
	// Jumio’s reference number for each scan. Max. lenght 36.
	ScanReference string `json:"scanReference"`
	Status        Status `json:"status"`
}
