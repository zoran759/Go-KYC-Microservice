package config_test

import (
	"modulus/kyc/main/config"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	assert := assert.New(t)

	err := config.FromFile("../kyc_dev.cfg")

	assert.NoError(err)

	err = config.FromFile("fake")

	assert.Error(err)
	assert.Equal("open fake: no such file or directory", err.Error())

	name := os.TempDir() + string(os.PathSeparator) + "tmp" + time.Now().Format("20060102150405")
	file, err := os.Create(name)

	assert.NoError(err)

	file.Close()
	err = config.FromFile(name)

	assert.Error(err)
	assert.Equal("empty "+name, err.Error())

	err = config.FromFile("file_test.go")

	assert.Error(err)
}
