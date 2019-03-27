package config_test

import (
	"testing"

	"modulus/kyc/common"
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

	updated, errs = config.Update(config.Config{
		config.ServiceSection: config.Options{
			"Port": "",
		},
	})

	assert.True(updated)
	assert.Empty(errs)

	port = config.ServicePort()

	assert.Equal("8080", port)
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
		testCase{
			name: "Option error 2",
			config: config.Config{
				"IDology": config.Options{
					"Host":             "host",
					"Username":         "username",
					"Password":         "password",
					"UseSummaryResult": "smile",
				},
			},
			updated: true,
			errs: []string{
				"IDology 'UseSummaryResult' option error: strconv.ParseBool: parsing \"smile\": invalid syntax",
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

func TestIsKnownName(t *testing.T) {
	type testCase struct {
		name     string
		testname string
		isknown  bool
	}

	testCases := []testCase{
		testCase{
			name:     "Known name",
			testname: string(common.Jumio),
			isknown:  true,
		},
		testCase{
			name:     "Unknown name",
			testname: "Fake",
		},
		testCase{
			name:     "Service section",
			testname: config.ServiceSection,
		},
		testCase{
			name:     "Not provider",
			testname: string(common.CipherTrace),
			isknown:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isknown := config.IsKnownName(tc.testname)
			assert.Equal(t, tc.isknown, isknown)
		})
	}
}

func TestGetOptions(t *testing.T) {
	assert := assert.New(t)

	cfg := config.Config{
		"IDology": config.Options{
			"Host":             "Host",
			"Username":         "Username",
			"Password":         "Password",
			"UseSummaryResult": "true",
		},
	}

	updated, err := config.Update(cfg)

	assert.True(updated)
	assert.Empty(err)

	type testCase struct {
		name    string
		section string
		options config.Options
	}

	testCases := []testCase{
		testCase{
			name:    "Valid options",
			section: "IDology",
			options: config.Options{
				"Host":             "Host",
				"Username":         "Username",
				"Password":         "Password",
				"UseSummaryResult": "true",
			},
		},
		testCase{
			name:    "Empty options",
			section: "Fake",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := config.GetOptions(tc.section)
			assert.Equal(tc.options, opts)
		})
	}
}

func TestGetConfig(t *testing.T) {
	assert := assert.New(t)

	sect := "IDology"
	cfg := config.Config{
		sect: config.Options{
			"Host":             "Host",
			"Username":         "Username",
			"Password":         "Password",
			"UseSummaryResult": "true",
		},
	}

	updated, err := config.Update(cfg)

	assert.True(updated)
	assert.Empty(err)

	testcfg := config.GetConfig()

	assert.Contains(testcfg, sect)

	opts := cfg[sect]
	testopts := testcfg[sect]

	assert.Equal(opts, testopts)
}
