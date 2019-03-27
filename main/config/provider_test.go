package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm"
	"modulus/kyc/integrations/complyadvantage"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/idology"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/synapsefi"
	"modulus/kyc/integrations/thomsonreuters"
	"modulus/kyc/integrations/trulioo"

	"github.com/stretchr/testify/assert"
)

var testValidConfig = `
[Coinfirm]
Host=https://api.coinfirm.io/v2
Email=info@example.com
Password=lbslgsbdgldfsnblfdnbhldf
Company=Foobar

[ComplyAdvantage]
Host=https://api.complyadvantage.com
APIkey=jkghjilgbsldbgsegbnskeflew
Fuzziness=0.3

[IdentityMind]
Host=https://staging.identitymind.com/im
Username=foobar
Password=fiweugbilewbgulwebglgblier

[IDology]
Host=https://web.idologylive.com/api/idiq.svc
Username=foobar
Password=gbisgubsigbs
UseSummaryResult=false

[Jumio]
BaseURL=https://lon.netverify.com/api/netverify/v2
Token=c278848a-0de3-35da-5be1-fa575c9b94c4
Secret=ifjsgsbdlgbudsglasbiuewgwivkgwequkkqww

[ShuftiPro]
Host=https://shuftipro.com/api/
ClientID=kdsgldghhlesrgnserlhnrhnerlhn
SecretKey=ihbgielrbglerablqwafblwebglerbgrle
CallbackURL=https://shuftipro.com/api

[Sum&Substance]
Host=https://test-api.sumsub.com
APIKey=GKTBNXNEPJHCXY

[SynapseFI]
Host=https://uat-api.synapsefi.com/v3.1/
ClientID=client_id_vfoewilavflaergaboebjqpblvrlevfrlbvlrbl
ClientSecret=client_secret_yufwufvwufvaefskufasvkfevwskf

[ThomsonReuters]
Host=https://rms-world-check-one-api-pilot.thomsonreuters.com/v1/
# Host=https://rms-world-check-one-api.thomsonreuters.com/v1/
APIkey=c278848a-0de3-35da-5be1-fa575c9b94c4
APIsecret=huefhgeikgbkabowqufikuywegr48735gt43ujvti3vtkqvwd7632qig5ruk32vqakf3r

[Trulioo]
Host=https://api.globaldatacompany.com
NAPILogin=foobar
NAPIPassword=67fuefyvqwurvweafkwi4r

[Config]
Port=8080

[CipherTrace]
URL=https://rest.ciphertrace.com
Key=jwgfydisavfikyagfweyo435g98iwukvyf43ifkuajsvdi723wvdkw
Username=superkey
`

func TestCreatePlatform(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "kyc_cfg_")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(testValidConfig))
	tmpfile.Close()

	Load(tmpfile.Name())

	coinfirmOpts := cfg.config[string(common.Coinfirm)]
	complyadvantageOpts := cfg.config[string(common.ComplyAdvantage)]
	identitymindOpts := cfg.config[string(common.IdentityMind)]
	idologyOpts := cfg.config[string(common.IDology)]
	jumioOpts := cfg.config[string(common.Jumio)]
	shuftiproOpts := cfg.config[string(common.ShuftiPro)]
	sumsubOpts := cfg.config[string(common.SumSub)]
	synapsefiOpts := cfg.config[string(common.SynapseFI)]
	thomsonreutersOpts := cfg.config[string(common.ThomsonReuters)]
	truliooOpts := cfg.config[string(common.Trulioo)]

	fuzziness, _ := strconv.ParseFloat(complyadvantageOpts["Fuzziness"], 32)
	useSummaryResult, _ := strconv.ParseBool(idologyOpts["UseSummaryResult"])

	type testCase struct {
		name     string
		provider common.KYCProvider
		platform common.KYCPlatform
		err      error
	}

	testCases := []testCase{
		testCase{
			name:     string(common.Coinfirm),
			provider: common.Coinfirm,
			platform: coinfirm.New(coinfirm.Config{
				Host:     coinfirmOpts["Host"],
				Email:    coinfirmOpts["Email"],
				Password: coinfirmOpts["Password"],
				Company:  coinfirmOpts["Company"],
			}),
		},
		testCase{
			name:     string(common.ComplyAdvantage),
			provider: common.ComplyAdvantage,
			platform: complyadvantage.New(complyadvantage.Config{
				Host:      complyadvantageOpts["Host"],
				APIkey:    complyadvantageOpts["APIkey"],
				Fuzziness: float32(fuzziness),
			}),
		},
		testCase{
			name:     string(common.IdentityMind),
			provider: common.IdentityMind,
			platform: identitymind.New(identitymind.Config{
				Host:     identitymindOpts["Host"],
				Username: identitymindOpts["Username"],
				Password: identitymindOpts["Password"],
			}),
		},
		testCase{
			name:     string(common.IDology),
			provider: common.IDology,
			platform: idology.New(idology.Config{
				Host:             idologyOpts["Host"],
				Username:         idologyOpts["Username"],
				Password:         idologyOpts["Password"],
				UseSummaryResult: useSummaryResult,
			}),
		},
		testCase{
			name:     string(common.Jumio),
			provider: common.Jumio,
			platform: jumio.New(jumio.Config{
				BaseURL: jumioOpts["BaseURL"],
				Token:   jumioOpts["Token"],
				Secret:  jumioOpts["Secret"],
			}),
		},
		testCase{
			name:     string(common.ShuftiPro),
			provider: common.ShuftiPro,
			platform: shuftipro.New(shuftipro.Config{
				Host:        shuftiproOpts["Host"],
				ClientID:    shuftiproOpts["ClientID"],
				SecretKey:   shuftiproOpts["SecretKey"],
				CallbackURL: shuftiproOpts["CallbackURL"],
			}),
		},
		testCase{
			name:     string(common.SumSub),
			provider: common.SumSub,
			platform: sumsub.New(sumsub.Config{
				Host:   sumsubOpts["Host"],
				APIKey: sumsubOpts["APIKey"],
			}),
		},
		testCase{
			name:     string(common.SynapseFI),
			provider: common.SynapseFI,
			platform: synapsefi.New(synapsefi.Config{
				Host:         synapsefiOpts["Host"],
				ClientID:     synapsefiOpts["ClientID"],
				ClientSecret: synapsefiOpts["ClientSecret"],
			}),
		},
		testCase{
			name:     string(common.ThomsonReuters),
			provider: common.ThomsonReuters,
			platform: thomsonreuters.New(thomsonreuters.Config{
				Host:      thomsonreutersOpts["Host"],
				APIkey:    thomsonreutersOpts["APIkey"],
				APIsecret: thomsonreutersOpts["APIsecret"],
			}),
		},
		testCase{
			name:     string(common.Trulioo),
			provider: common.Trulioo,
			platform: trulioo.New(trulioo.Config{
				Host:         truliooOpts["Host"],
				NAPILogin:    truliooOpts["NAPILogin"],
				NAPIPassword: truliooOpts["NAPIPassword"],
			}),
		},
		testCase{
			name:     "Invalid provider",
			provider: common.KYCProvider("Fake"),
			err:      errors.New("Fake is missing configuration validation"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := createPlatform(tc.provider)
			assert.Equal(t, tc.platform, p)
			assert.Equal(t, tc.err, err)
		})
	}
}
