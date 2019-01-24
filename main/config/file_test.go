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

	cfg, err := config.FromFile("../kyc_dev.cfg")

	assert.NoError(err)
	assert.NotNil(cfg)

	cfg, err = config.FromFile("fake")

	assert.Error(err)
	assert.Equal("open fake: no such file or directory", err.Error())
	assert.Nil(cfg)

	tmpfile, err := ioutil.TempFile("", "kyc")
	assert.NoError(err)

	defer os.Remove(tmpfile.Name())

	err = tmpfile.Close()
	assert.NoError(err)

	cfg, err = config.FromFile(tmpfile.Name())

	assert.Error(err)
	assert.Equal("empty "+tmpfile.Name(), err.Error())
	assert.Nil(cfg)

	cfg, err = config.FromFile("file_test.go")

	assert.Error(err)
	assert.Equal("parsing failed at line 1 'package config_test': not proper config string", err.Error())
	assert.Nil(cfg)
}
