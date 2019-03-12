package common

import "time"

// Time defines the model for time values.
type Time time.Time

// MarshalJSON implements the json.Marshaler interface for the common.Time type.
func (t Time) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface for the common.Time type.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		return nil
	}
	return (*time.Time)(t).UnmarshalJSON(data)
}

// The returned string is meant for debugging; for a stable serialized
// representation, use t.Format with an explicit format string.
func (t Time) String() string {
	return time.Time(t).String()
}

// Format returns time string formatted according to the provided layout.
func (t Time) Format(layout string) string {
	if !time.Time(t).IsZero() {
		return time.Time(t).Format(layout)
	}
	return ""
}
