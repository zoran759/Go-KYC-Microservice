package jumio

import (
	"encoding/base64"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Describe("New", func() {
		Specify("should properly create service object", func() {
			config := Config{
				BaseURL: "fake_baseURL",
				Token:   "fake_token",
				Secret:  "fake_secret",
			}

			s := &service{
				baseURL:     config.BaseURL,
				credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
			}

			cc := New(config)
			ts := cc.(*service)

			Expect(s).To(Equal(ts))
		})
	})
})
