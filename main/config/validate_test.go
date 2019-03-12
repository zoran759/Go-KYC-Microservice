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
		"CallbackURL": "host",
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
	err := validate(validConfig)

	assert.NoError(t, err)
}

func TestVerifyComplyAdvantage(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		string(common.ComplyAdvantage): Options{
			"APIkey":    "key",
			"Fuzziness": "0",
		},
	}

	err := validate(config)

	assert.Error(err)
	assert.Equal(reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(`ComplyAdvantage configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.ComplyAdvantage): Options{
			"Host":      "host",
			"Fuzziness": "0",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`ComplyAdvantage configuration error: missing or empty option 'APIkey'`, err.Error())

	config = Config{
		string(common.ComplyAdvantage): Options{
			"Host":   "host",
			"APIkey": "key",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`ComplyAdvantage configuration error: missing or empty option 'Fuzziness'`, err.Error())
}

func TestVerifyIdentityMind(t *testing.T) {
	config := Config{
		string(common.IdentityMind): Options{
			"Username": "fakeuser",
			"Password": "fakepassword",
		},
	}

	err := validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(t, `IdentityMind configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.IdentityMind): Options{
			"Host":     "host",
			"Password": "fakepassword",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `IdentityMind configuration error: missing or empty option 'Username'`, err.Error())

	config = Config{
		string(common.IdentityMind): Options{
			"Host":     "host",
			"Username": "fakeuser",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `IdentityMind configuration error: missing or empty option 'Password'`, err.Error())
}

func TestVerifyIDology(t *testing.T) {
	config := Config{
		string(common.IDology): Options{
			"Username":         "fakeuser",
			"Password":         "fakepassword",
			"UseSummaryResult": "false",
		},
	}

	err := validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(t, `IDology configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.IDology): Options{
			"Host":             "host",
			"Password":         "fakepassword",
			"UseSummaryResult": "false",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `IDology configuration error: missing or empty option 'Username'`, err.Error())

	config = Config{
		string(common.IDology): Options{
			"Host":             "host",
			"Username":         "fakeuser",
			"UseSummaryResult": "false",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `IDology configuration error: missing or empty option 'Password'`, err.Error())

	config = Config{
		string(common.IDology): Options{
			"Host":     "host",
			"Username": "fakeuser",
			"Password": "fakepassword",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `IDology configuration error: missing or empty option 'UseSummaryResult'`, err.Error())
}

func TestVerifyJumio(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		string(common.Jumio): Options{
			"Token":  "token",
			"Secret": "secret",
		},
	}

	err := validate(config)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(`Jumio configuration error: missing or empty option 'BaseURL'`, err.Error())

	config = Config{
		string(common.Jumio): Options{
			"BaseURL": "base_url",
			"Secret":  "secret",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`Jumio configuration error: missing or empty option 'Token'`, err.Error())

	config = Config{
		string(common.Jumio): Options{
			"BaseURL": "base_url",
			"Token":   "token",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`Jumio configuration error: missing or empty option 'Secret'`, err.Error())
}

func TestVerifyShuftiPro(t *testing.T) {
	config := Config{
		string(common.ShuftiPro): Options{
			"ClientID":    "fakeid",
			"SecretKey":   "fakekey",
			"RedirectURL": "host",
		},
	}

	err := validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(t, `ShuftiPro configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.ShuftiPro): Options{
			"Host":        "host",
			"SecretKey":   "fakekey",
			"RedirectURL": "host",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `ShuftiPro configuration error: missing or empty option 'ClientID'`, err.Error())

	config = Config{
		string(common.ShuftiPro): Options{
			"Host":        "host",
			"ClientID":    "fakeid",
			"RedirectURL": "host",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `ShuftiPro configuration error: missing or empty option 'SecretKey'`, err.Error())

	config = Config{
		string(common.ShuftiPro): Options{
			"Host":      "host",
			"ClientID":  "fakeid",
			"SecretKey": "fakekey",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `ShuftiPro configuration error: missing or empty option 'CallbackURL'`, err.Error())
}

func TestVerifySumSub(t *testing.T) {
	config := Config{
		string(common.SumSub): Options{
			"APIKey": "fakekey",
		},
	}

	err := validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(t, `Sum&Substance configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.SumSub): Options{
			"Host": "host",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `Sum&Substance configuration error: missing or empty option 'APIKey'`, err.Error())
}

func TestVerifySynapseFI(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		string(common.SynapseFI): Options{
			"ClientID":     "clientID",
			"ClientSecret": "secret",
		},
	}

	err := validate(config)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(`SynapseFI configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.SynapseFI): Options{
			"Host":         "host",
			"ClientSecret": "secret",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`SynapseFI configuration error: missing or empty option 'ClientID'`, err.Error())

	config = Config{
		string(common.SynapseFI): Options{
			"Host":     "host",
			"ClientID": "clientID",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`SynapseFI configuration error: missing or empty option 'ClientSecret'`, err.Error())
}

func TestVerifyThomsonReuters(t *testing.T) {
	assert := assert.New(t)

	config := Config{
		string(common.ThomsonReuters): Options{
			"APIkey":    "key",
			"APIsecret": "secret",
		},
	}

	err := validate(config)
	assert.Error(err)
	assert.Equal(reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(`ThomsonReuters configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.ThomsonReuters): Options{
			"Host":      "host",
			"APIsecret": "secret",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`ThomsonReuters configuration error: missing or empty option 'APIkey'`, err.Error())

	config = Config{
		string(common.ThomsonReuters): Options{
			"Host":   "host",
			"APIkey": "key",
		},
	}

	err = validate(config)
	assert.Error(err)
	assert.Equal(`ThomsonReuters configuration error: missing or empty option 'APIsecret'`, err.Error())
}

func TestVerifyTrulioo(t *testing.T) {
	config := Config{
		string(common.Trulioo): Options{
			"NAPILogin":    "fakelogin",
			"NAPIPassword": "fakepassword",
		},
	}

	err := validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(ErrMissingOption{}), reflect.TypeOf(err))
	assert.Equal(t, `Trulioo configuration error: missing or empty option 'Host'`, err.Error())

	config = Config{
		string(common.Trulioo): Options{
			"Host":         "host",
			"NAPIPassword": "fakepassword",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `Trulioo configuration error: missing or empty option 'NAPILogin'`, err.Error())

	config = Config{
		string(common.Trulioo): Options{
			"Host":      "host",
			"NAPILogin": "fakelogin",
		},
	}

	err = validate(config)
	assert.NotNil(t, err)
	assert.Equal(t, `Trulioo configuration error: missing or empty option 'NAPIPassword'`, err.Error())
}
