package config_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"modulus/kyc/main/config"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {

	fakefile := "../fake"
	zerofile := "../empty"
	hugefile := "../huge"
	invalidfile := "../invalid"
	cfgerrorsfile := "../witherrs"
	upderrsfile := "../upderrs"

	f, err := os.Create(zerofile)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	f, err = os.Create(hugefile)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(strings.Repeat("s", 1<<20+1))
	f.Close()

	f, err = os.Create(invalidfile)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("# Invalid config file content\n")
	f.Close()

	f, err = os.Create(cfgerrorsfile)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("# Config file with errors\nxyz")
	f.Close()

	f, err = os.Create(upderrsfile)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("# Wrong option\n[Coinfirm]\nFakeOpt=value\n")
	f.Close()

	defer os.Remove(fakefile)
	defer os.Remove(zerofile)
	defer os.Remove(hugefile)
	defer os.Remove(invalidfile)
	defer os.Remove(cfgerrorsfile)
	defer os.Remove(upderrsfile)

	log.SetFlags(log.Ldate)

	type testCase struct {
		name     string
		filename string
		logger   bytes.Buffer
		logs     []byte
	}

	testCases := []testCase{
		testCase{
			name:     "Valid file",
			filename: "../" + config.DefaultDevFile,
		},
		testCase{
			name:     "Invalid file",
			filename: invalidfile,
			logs:     []byte(time.Now().Format("2006/01/02") + " config doesn't contains a proper configuration\n"),
		},
		testCase{
			name:     "Config with errors",
			filename: cfgerrorsfile,
			logs: []byte(time.Now().Format("2006/01/02") + " [config parser] error at line 2: 'xyz'\n" +
				time.Now().Format("2006/01/02") + " config doesn't contains a proper configuration\n"),
		},
		testCase{
			name:     "Wrong option",
			filename: upderrsfile,
			logs:     []byte(time.Now().Format("2006/01/02") + " Coinfirm: unknown option 'FakeOpt' was filtered out\n"),
		},
		testCase{
			name:     "File not exists",
			filename: fakefile,
			logs:     []byte(time.Now().Format("2006/01/02") + " INFO: '../fake' empty configuration file created. No config loaded\n"),
		},
		testCase{
			name:     "Inaccessible directory",
			filename: "FAKE/fake",
			logs:     []byte(time.Now().Format("2006/01/02") + " open FAKE/fake: no such file or directory\n"),
		},
		testCase{
			name:     "Inaccessible file",
			filename: "/root/test",
			logs:     []byte(time.Now().Format("2006/01/02") + " stat /root/test: permission denied\n"),
		},
		testCase{
			name:     "Empty file",
			filename: zerofile,
			logs:     []byte(time.Now().Format("2006/01/02") + " WARNING! '../empty' configuration file is empty. No config loaded\n"),
		},
		testCase{
			name:     "Huge file",
			filename: hugefile,
			logs:     []byte(time.Now().Format("2006/01/02") + " WARNING! '../huge' size 1048577 exceeded the limit of 1 Mb. No config loaded\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.SetOutput(&tc.logger)
			config.Load(tc.filename)
			assert.Equal(t, tc.logs, tc.logger.Bytes())
		})
	}
}
