package model

// List of ProviderSourceStatus values.
const (
	ActiveSource  ProviderSourceStatus = "ACTIVE"
	DeletedSource ProviderSourceStatus = "DELETED"
	HiddenSource  ProviderSourceStatus = "HIDDEN"
)

// List of SubscriptionCategory values.
const (
	Premium  SubscriptionCategory = "PREMIUM"
	Standard SubscriptionCategory = "STANDARD"
)

// ProviderSourceStatus represents the provider source status enumeration.
type ProviderSourceStatus string

// SubscriptionCategory represents the subscription category type enumeration.
type SubscriptionCategory string

// ProviderSourceTypeCategoryDetail represents a provider source type category detail.
type ProviderSourceTypeCategoryDetail struct {
	ID                  string               `json:"identifier"`
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	ProviderSourceTypes []ProviderSourceType `json:"providerSourceTypes"`
}

// ProviderSourceType represents a provider source type.
type ProviderSourceType struct {
	ID       string                           `json:"identifier"`
	Name     string                           `json:"name"`
	Category ProviderSourceTypeCategoryDetail `json:"category"`
}

// Provider represents a provider for a source.
type Provider struct {
	ID     string `json:"identifier"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Master bool   `json:"master"`
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
	CreationDate         string               `json:"creationDate"`
	SubscriptionCategory SubscriptionCategory `json:"subscriptionCategory"`
}

// ProviderDetail represents a provider.
type ProviderDetail struct {
	ID      string           `json:"identifier"`
	Code    string           `json:"code"`
	Name    string           `json:"name"`
	Master  bool             `json:"master"`
	Sources []ProviderSource `json:"sources"`
}

// ProviderDetails represents a list of available providers and their sources.
type ProviderDetails []ProviderDetail
