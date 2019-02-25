package config

import (
	"reflect"
	"testing"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

var validConfig = Config{
	string(common.IdentityMind): Options{
		"Host":     "host",
		"Username": "fakeuser",
		"Password": "fakepassword",
	},
	string(common.IDology): Options{
		"Host":             "host",
		"Username":         "fakeuser",
		"Password":         "fakepassword",
		"UseSummaryResult": "false",
	},
	string(common.ShuftiPro): Options{
		"Host":        "host",
		"ClientID":    "fakeid",
		"SecretKey":   "fakekey",
		"RedirectURL": "host",
	},
	string(common.SumSub): Options{
		"Host":   "host",
		"APIKey": "fakekey",
	},
	string(common.Trulioo): Options{
		"Host":         "host",
		"NAPILogin":    "fakelogin",
		"NAPIPassword": "fakepassword",
	},
}

func TestVerifySuccess(t *testing.T) {
	assert := assert.New(t)

	for p, oo := range validConfig {
		err := validateProvider(common.KYCProvider(p), oo)
		assert.NoError(err)
	}
}

func TestVerifyComplyAdvantage(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"APIkey":    "key",
		"Fuzziness": "0",
	}

	err := validateProvider(common.ComplyAdvantage, opts)

	assert.Error(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the ComplyAdvantage provider`, err.Error())

	opts = Options{
		"Host":      "host",
		"Fuzziness": "0",
	}

	err = validateProvider(common.ComplyAdvantage, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'APIkey' for the ComplyAdvantage provider`, err.Error())

	opts = Options{
		"Host":   "host",
		"APIkey": "key",
	}

	err = validateProvider(common.ComplyAdvantage, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'Fuzziness' for the ComplyAdvantage provider`, err.Error())
}

func TestVerifyIdentityMind(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"Username": "fakeuser",
		"Password": "fakepassword",
	}

	err := validateProvider(common.IdentityMind, opts)
	assert.NotNil(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the IdentityMind provider`, err.Error())

	opts = Options{
		"Host":     "host",
		"Password": "fakepassword",
	}

	err = validateProvider(common.IdentityMind, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'Username' for the IdentityMind provider`, err.Error())

	opts = Options{
		"Host":     "host",
		"Username": "fakeuser",
	}

	err = validateProvider(common.IdentityMind, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'Password' for the IdentityMind provider`, err.Error())
}

func TestVerifyIDology(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"Username":         "fakeuser",
		"Password":         "fakepassword",
		"UseSummaryResult": "false",
	}

	err := validateProvider(common.IDology, opts)
	assert.NotNil(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the IDology provider`, err.Error())

	opts = Options{
		"Host":             "host",
		"Password":         "fakepassword",
		"UseSummaryResult": "false",
	}

	err = validateProvider(common.IDology, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'Username' for the IDology provider`, err.Error())

	opts = Options{
		"Host":             "host",
		"Username":         "fakeuser",
		"UseSummaryResult": "false",
	}

	err = validateProvider(common.IDology, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'Password' for the IDology provider`, err.Error())

	opts = Options{
		"Host":     "host",
		"Username": "fakeuser",
		"Password": "fakepassword",
	}

	err = validateProvider(common.IDology, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'UseSummaryResult' for the IDology provider`, err.Error())
}

func TestVerifyJumio(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"Token":  "token",
		"Secret": "secret",
	}

	err := validateProvider(common.Jumio, opts)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'BaseURL' for the Jumio provider`, err.Error())

	opts = Options{
		"BaseURL": "base_url",
		"Secret":  "secret",
	}

	err = validateProvider(common.Jumio, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'Token' for the Jumio provider`, err.Error())

	opts = Options{
		"BaseURL": "base_url",
		"Token":   "token",
	}

	err = validateProvider(common.Jumio, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'Secret' for the Jumio provider`, err.Error())
}

func TestVerifyShuftiPro(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"ClientID":    "fakeid",
		"SecretKey":   "fakekey",
		"RedirectURL": "host",
	}

	err := validateProvider(common.ShuftiPro, opts)
	assert.NotNil(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the ShuftiPro provider`, err.Error())

	opts = Options{
		"Host":        "host",
		"SecretKey":   "fakekey",
		"RedirectURL": "host",
	}

	err = validateProvider(common.ShuftiPro, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'ClientID' for the ShuftiPro provider`, err.Error())

	opts = Options{
		"Host":        "host",
		"ClientID":    "fakeid",
		"RedirectURL": "host",
	}

	err = validateProvider(common.ShuftiPro, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'SecretKey' for the ShuftiPro provider`, err.Error())

	opts = Options{
		"Host":      "host",
		"ClientID":  "fakeid",
		"SecretKey": "fakekey",
	}

	err = validateProvider(common.ShuftiPro, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'RedirectURL' for the ShuftiPro provider`, err.Error())
}

func TestVerifySumSub(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"APIKey": "fakekey",
	}

	err := validateProvider(common.SumSub, opts)
	assert.NotNil(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the Sum&Substance provider`, err.Error())

	opts = Options{
		"Host": "host",
	}

	err = validateProvider(common.SumSub, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'APIKey' for the Sum&Substance provider`, err.Error())
}

func TestVerifySynapseFI(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"ClientID":     "clientID",
		"ClientSecret": "secret",
	}

	err := validateProvider(common.SynapseFI, opts)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the SynapseFI provider`, err.Error())

	opts = Options{
		"Host":         "host",
		"ClientSecret": "secret",
	}

	err = validateProvider(common.SynapseFI, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'ClientID' for the SynapseFI provider`, err.Error())

	opts = Options{
		"Host":     "host",
		"ClientID": "clientID",
	}

	err = validateProvider(common.SynapseFI, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'ClientSecret' for the SynapseFI provider`, err.Error())
}

func TestVerifyThomsonReuters(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"APIkey":    "key",
		"APIsecret": "secret",
	}

	err := validateProvider(common.ThomsonReuters, opts)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the ThomsonReuters provider`, err.Error())

	opts = Options{
		"Host":      "host",
		"APIsecret": "secret",
	}

	err = validateProvider(common.ThomsonReuters, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'APIkey' for the ThomsonReuters provider`, err.Error())

	opts = Options{
		"Host":   "host",
		"APIkey": "key",
	}

	err = validateProvider(common.ThomsonReuters, opts)
	assert.Error(err)
	assert.Equal(`missing or empty option 'APIsecret' for the ThomsonReuters provider`, err.Error())
}

func TestVerifyTrulioo(t *testing.T) {
	assert := assert.New(t)

	opts := Options{
		"NAPILogin":    "fakelogin",
		"NAPIPassword": "fakepassword",
	}

	err := validateProvider(common.Trulioo, opts)
	assert.NotNil(err)
	assert.Equal(reflect.TypeOf(MissingOptionError{}), reflect.TypeOf(err))
	assert.Equal(`missing or empty option 'Host' for the Trulioo provider`, err.Error())

	opts = Options{
		"Host":         "host",
		"NAPIPassword": "fakepassword",
	}

	err = validateProvider(common.Trulioo, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'NAPILogin' for the Trulioo provider`, err.Error())

	opts = Options{
		"Host":      "host",
		"NAPILogin": "fakelogin",
	}

	err = validateProvider(common.Trulioo, opts)
	assert.NotNil(err)
	assert.Equal(`missing or empty option 'NAPIPassword' for the Trulioo provider`, err.Error())
}
