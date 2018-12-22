package model

// List of ProviderSourceStatus values.
const (
	ActiveSource  ProviderSourceStatus = "ACTIVE"
	DeletedSource ProviderSourceStatus = "DELETED"
	HiddenSource  ProviderSourceStatus = "HIDDEN"
)

// List of ProviderType values.
const (
	WatchList       ProviderType = "WATCHLIST"
	PassportCheck   ProviderType = "PASSPORT_CHECK"
	ClientWatchList ProviderType = "CLIENT_WATCHLIST"
)

// List of SubscriptionCategory values.
const (
	Premium  SubscriptionCategory = "PREMIUM"
	Standard SubscriptionCategory = "STANDARD"
)

// ProviderDetails represents a list of available providers and their sources.
type ProviderDetails []ProviderDetail

// ProviderDetail represents a provider.
type ProviderDetail struct {
	ID      string           `json:"identifier"`
	Code    string           `json:"code"`
	Name    string           `json:"name"`
	Master  bool             `json:"master"`
	Sources []ProviderSource `json:"sources"`
}

// ProviderSource represents a source that belongs to a provider.
type ProviderSource struct {
	ID                   string               `json:"identifier"`
	ImportID             string               `json:"importIdentifier"`
	Name                 string               `json:"name"`
	Abbreviation         string               `json:"abbreviation"`
	Provider             Provider             `json:"provider"`
	Type                 ProviderSourceType   `json:"type"`
	ProviderSourceStatus ProviderSourceStatus `json:"providerSourceStatus"`
	RegionOfAuthority    string               `json:"regionOfAuthority"`
	SubscriptionCategory SubscriptionCategory `json:"subscriptionCategory"`
	CreationDate         string               `json:"creationDate"`
}

// Provider represents a provider for a source.
type Provider struct {
	ID     string `json:"identifier"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Master bool   `json:"master"`
}

// ProviderSourceStatus represents the provider source status enumeration.
type ProviderSourceStatus string

// SubscriptionCategory represents the subscription category type enumeration.
type SubscriptionCategory string

// ProviderSourceType represents a provider source type.
type ProviderSourceType struct {
	ID       string                           `json:"identifier"`
	Name     string                           `json:"name"`
	Category ProviderSourceTypeCategoryDetail `json:"category"`
}

// ProviderSourceTypeCategoryDetail represents a provider source type category detail.
type ProviderSourceTypeCategoryDetail struct {
	ID                  string               `json:"identifier"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	ProviderSourceTypes []ProviderSourceType `json:"providerSourceTypes"`
}

// ProviderType represents the provider type enumeration.
type ProviderType string
