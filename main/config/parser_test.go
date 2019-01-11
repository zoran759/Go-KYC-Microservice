package config

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawConfig = `
# This is a sample of the valid config.
#

[ComplyAdvantage]
Host=https://api.complyadvantage.com
APIkey=bBF7TsnK2xCIMMNNbgUNcDwvSRSIbLT9
Fuzziness=0.3

[IdentityMind]
# This is the sandbox url.
# You should change it to the production url when deploying to production servers.
Host=https://sandbox.identitymind.com/im
Username=modulusglobal
Password=64117e699462ce859d970648461a625bc6a6f3cb

[IDology]
Host=https://web.idologylive.com/api/idiq.svc
Username=modulus.dev2
Password=}$tRPfT1sZQmU@uh8@
UseSummaryResult=false`

var rawConfigWithEmptyName = `
[IdentityMind]
# This is the sandbox url.
# You should change it to the production url when deploying to production servers.
Host=https://sandbox.identitymind.com/im
Username=modulusglobal
Password=64117e699462ce859d970648461a625bc6a6f3cb

[]
Host=https://web.idologylive.com/api/idiq.svc
Username=modulus.dev2
Password=}$tRPfT1sZQmU@uh8@
UseSummaryResult=false`

var rawConfigWithUnknownName = `
[IdentityMind]
# This is the sandbox url.
# You should change it to the production url when deploying to production servers.
Host=https://sandbox.identitymind.com/im
Username=modulusglobal
Password=64117e699462ce859d970648461a625bc6a6f3cb

[Foobar]
Host=https://web.idologylive.com/api/idiq.svc
Username=modulus.dev2
Password=}$tRPfT1sZQmU@uh8@
UseSummaryResult=false`

var rawConfigWithStandaloneOption = `
# This is an example of the config with standalone options.

# This is the sandbox url.
# You should change it to the production url when deploying to production servers.
Host=https://sandbox.identitymind.com/im
Username=modulusglobal
Password=64117e699462ce859d970648461a625bc6a6f3cb

[Foobar]
Host=https://web.idologylive.com/api/idiq.svc
Username=modulus.dev2
Password=}$tRPfT1sZQmU@uh8@
UseSummaryResult=false`

var rawConfigWithWrongNameFormat = `
[IdentityMind
# This is the sandbox url.
# You should change it to the production url when deploying to production servers.
Host=https://sandbox.identitymind.com/im
Username=modulusglobal
Password=64117e699462ce859d970648461a625bc6a6f3cb`

var emptyRawConfig = ``

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)

	reader := strings.NewReader(rawConfig)

	cfg, err := parseConfig(reader)

	assert.NoError(err)
	assert.NotNil(cfg)

	reader = strings.NewReader(rawConfigWithEmptyName)

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 9 '[]': missing KYC provider name in the config", err.Error())
	assert.Nil(cfg)

	reader = strings.NewReader(rawConfigWithUnknownName)

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 9 '[Foobar]': unknown KYC provider name in the config", err.Error())
	assert.Nil(cfg)

	reader = strings.NewReader(rawConfigWithStandaloneOption)

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 6 'Host=https://sandbox.identitymind.com/im': standalone option string", err.Error())
	assert.Nil(cfg)

	reader = strings.NewReader(rawConfigWithWrongNameFormat)

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 2 '[IdentityMind': not proper config string", err.Error())
	assert.Nil(cfg)

	reader = strings.NewReader(emptyRawConfig)

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 0 '': config is empty", err.Error())
	assert.Nil(cfg)

	cfg, err = parseConfig(nil)

	assert.Error(err)
	assert.Equal("the config source is nil", err.Error())
	assert.Nil(cfg)

	b := make([]rune, bufio.MaxScanTokenSize)
	reader = strings.NewReader(string(b))

	cfg, err = parseConfig(reader)

	assert.Error(err)
	assert.Equal("parsing failed at line 0 '': bufio.Scanner: token too long", err.Error())
	assert.Nil(cfg)
}
