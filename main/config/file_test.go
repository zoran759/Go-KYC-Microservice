package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"modulus/kyc/main/config"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	assert := assert.New(t)

	err := config.FromFile("../kyc_dev.cfg")

	assert.NoError(err)
	assert.NotEmpty(config.Cfg)

	err = config.FromFile("fake")

	assert.Error(err)
	assert.Equal("open fake: no such file or directory", err.Error())

	tmpfile, err := ioutil.TempFile("", "kyc")
	assert.NoError(err)

	defer os.Remove(tmpfile.Name())

	err = tmpfile.Close()
	assert.NoError(err)

	err = config.FromFile(tmpfile.Name())

	assert.Error(err)
	assert.Equal("empty "+tmpfile.Name(), err.Error())

	err = config.FromFile("file_test.go")

	assert.Error(err)
	assert.Equal("parsing failed at line 1 'package config_test': not proper config string", err.Error())
}
