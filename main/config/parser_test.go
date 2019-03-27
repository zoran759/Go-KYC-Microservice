package config

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validCfg = `# Test valid configuration.
[Coinfirm]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

[ComplyAdvantage]
Host=https://api.complyadvantage.com
APIkey=G7khJJGD7kHHLnsflLHLNlvfslng6LLsf
Fuzziness=0.5

#Some empty lines.


# End.`

var orphanOptsCfg = `# Test configuration with orphaned options without a section header.
Host=https://web.idologylive.com/api/idiq.svc
Username=johndoe
Password=G87jkb$k@lljfgUHB8c
UseSummaryResult=false

[Sum&Substance]
Host=https://test-api.sumsub.com
APIKey=MKDSFXLNXLBNXL

# End.`

var emptySectCfg = `# Test configuration with empty section name.
[]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

[Sum&Substance]
Host=https://test-api.sumsub.com
APIKey=MKDSFXLNXLBNXL

# End.`

var invalidSectCfg = `# Test configuration with invalid section name.
[Foobar]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

[Sum&Substance]
Host=https://test-api.sumsub.com
APIKey=MKDSFXLNXLBNXL

[     ]
BaseURL=https://lon.netverify.com/api/netverify/v2
Token=ea454c7a-013f-440c-7fea-d1f05c98df97
Secret=YUOdnflKGFDLNlsnldLFIHLFNslewtHJj

# End.`

var errorInNameCfg = `# Test configuration with error in section name.
[Sum&Substance]
Host=https://test-api.sumsub.com
APIKey=MKDSFXLNXLBNXL

# This simulates forgotten closing bracket.
# The section and its belonging options will be omitted.
[Coinfirm
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

# End.`

var errorInCfg = `# Test configuration with errors.
[Sum&Substance]
Host=https://test-api.sumsub.com

# The following line contains format error.
APIKey MKDSFXLNXLBNXL

# The following line contains occasionally typed characters.
jhsfja

[Coinfirm]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

# End.`

var malformedCfg = `# Test malformed configuration.
Coinfirm]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Foobar

[ComplyAdvantage
Host=https://api.complyadvantage.com
APIkey=G7khJJGD7kHHLnsflLHLNlvfslng6LLsf
Fuzziness=0.5`

var scanErrCfg = `
[Coinfirm]
Host=https://api.coinfirm.io/v2
Email=info@foobar.com
Password=Aeb5@8c3b1d8b4633bcb85224e822609
Company=Fooba` + strings.Repeat("r", 64*1024)

func TestParseConfig(t *testing.T) {
	type testCase struct {
		name string
		src  io.Reader
		cfg  Config
		err  error
	}

	testCases := []testCase{
		testCase{
			name: "Valid config",
			src:  strings.NewReader(validCfg),
			cfg: Config{
				"Coinfirm": Options{
					"Host":     "https://api.coinfirm.io/v2",
					"Email":    "info@foobar.com",
					"Password": "Aeb5@8c3b1d8b4633bcb85224e822609",
					"Company":  "Foobar",
				},
				"ComplyAdvantage": Options{
					"Host":      "https://api.complyadvantage.com",
					"APIkey":    "G7khJJGD7kHHLnsflLHLNlvfslng6LLsf",
					"Fuzziness": "0.5",
				},
			},
		},
		testCase{
			name: "Missing source",
			err:  errors.New("the config source is nil"),
		},
		testCase{
			name: "Orphaned options",
			src:  strings.NewReader(orphanOptsCfg),
			cfg: Config{
				"Sum&Substance": Options{
					"Host":   "https://test-api.sumsub.com",
					"APIKey": "MKDSFXLNXLBNXL",
				},
			},
		},
		testCase{
			name: "Empty section",
			src:  strings.NewReader(emptySectCfg),
			cfg: Config{
				"Sum&Substance": Options{
					"Host":   "https://test-api.sumsub.com",
					"APIKey": "MKDSFXLNXLBNXL",
				},
			},
		},
		testCase{
			name: "Invalid section",
			src:  strings.NewReader(invalidSectCfg),
			cfg: Config{
				"Sum&Substance": Options{
					"Host":   "https://test-api.sumsub.com",
					"APIKey": "MKDSFXLNXLBNXL",
				},
			},
		},
		testCase{
			name: "Error in name",
			src:  strings.NewReader(errorInNameCfg),
			cfg: Config{
				"Sum&Substance": Options{
					"Host":   "https://test-api.sumsub.com",
					"APIKey": "MKDSFXLNXLBNXL",
				},
			},
		},
		testCase{
			name: "Errors in config",
			src:  strings.NewReader(errorInCfg),
			cfg: Config{
				"Sum&Substance": Options{
					"Host": "https://test-api.sumsub.com",
				},
				"Coinfirm": Options{
					"Host":     "https://api.coinfirm.io/v2",
					"Email":    "info@foobar.com",
					"Password": "Aeb5@8c3b1d8b4633bcb85224e822609",
					"Company":  "Foobar",
				},
			},
		},
		testCase{
			name: "Malformed config",
			src:  strings.NewReader(malformedCfg),
			err:  errors.New("config doesn't contains a proper configuration"),
		},
		testCase{
			name: "Scan error",
			src:  strings.NewReader(scanErrCfg),
			err: ParseError{
				strnum: 5,
				err:    "bufio.Scanner: token too long",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := parseConfig(tc.src)
			assert.Equal(t, tc.cfg, cfg)
			assert.Equal(t, tc.err, err)
		})
	}
}
