package config_test

import (
	"testing"

	"modulus/kyc/main/config"

	"github.com/stretchr/testify/assert"
)

func TestServicePort(t *testing.T) {
	assert := assert.New(t)

	port := config.ServicePort()

	assert.Equal(config.DefaultPort, port)

	updated, errs := config.Update(config.Config{
		config.ServiceSection: config.Options{
			"Port": "8999",
		},
	})

	assert.True(updated)
	assert.Empty(errs)

	port = config.ServicePort()

	assert.Equal("8999", port)
}

func TestUpdate(t *testing.T) {
	type testCase struct {
		name    string
		config  config.Config
		updated bool
		errs    []string
	}

	testCases := []testCase{
		testCase{
			name:   "Empty config",
			config: config.Config{},
			errs: []string{
				"no config update data provided",
			},
		},
		testCase{
			name: "Unknown section",
			config: config.Config{
				"Foobar": config.Options{
					"key1": "val1",
					"key2": "val2",
				},
			},
			errs: []string{
				"unknown config section: Foobar",
			},
		},
		testCase{
			name: "Empty section",
			config: config.Config{
				"IDology": config.Options{},
			},
			errs: []string{
				"empty config section: IDology",
			},
		},
		testCase{
			name: "Filtered out option",
			config: config.Config{
				"ComplyAdvantage": config.Options{
					"Host":      "host",
					"APIkey":    "key",
					"Fuzziness": "0",
					"Unknown":   "option",
				},
			},
			updated: true,
			errs: []string{
				"ComplyAdvantage: unknown option 'Unknown' was filtered out",
			},
		},
		testCase{
			name: "Empty config after filtering",
			config: config.Config{
				"ComplyAdvantage": config.Options{
					"Unknown": "option",
				},
			},
			errs: []string{
				"ComplyAdvantage: unknown option 'Unknown' was filtered out",
			},
		},
		testCase{
			name: "License update",
			config: config.Config{
				"Config": config.Options{
					"License": "invalid",
				},
			},
			updated: true,
			errs: []string{
				"during license validation: LicenseKey is not registered",
			},
		},
		testCase{
			name: "Option error",
			config: config.Config{
				"ComplyAdvantage": config.Options{
					"Host":      "host",
					"APIkey":    "key",
					"Fuzziness": "0,5",
				},
			},
			updated: true,
			errs: []string{
				"ComplyAdvantage 'Fuzziness' option error: strconv.ParseFloat: parsing \"0,5\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updated, errs := config.Update(tc.config)
			assert.Equal(t, tc.updated, updated)
			assert.Equal(t, tc.errs, errs)
		})
	}
}
