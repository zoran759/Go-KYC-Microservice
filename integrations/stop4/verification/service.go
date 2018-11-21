package verification

import (
	"encoding/json"
	"modulus/kyc/http"
	"github.com/gofrs/uuid"
	"net/url"
	"reflect"
	"regexp"
)

type service struct {
	config Config
}

// NewService constructs a new verification service object.
func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

func (service service) Verify(request RegistrationRequest) (response *Response, err error) {

	fields := make(map[string]string)
	convertModelToForm(request.CustomerInformation, fields)

	service.AttachAuthPartToRequest(&request)
	request.UserNumber, err = generateUserNumber()
	if err != nil {
		return nil, err
	}

	host := service.config.Host + "/customerregistration"

	form := url.Values{}
	form.Add("merchant_id", request.MerchantID)
	form.Add("password", request.Password)
	form.Add("reg_ip_address", request.RegIPAddress)
	form.Add("reg_date", request.RegDate)
	form.Add("user_number", request.UserNumber.String())
	for key, val := range fields {
		form.Add(key, val)
	}

	responseStatus, responseBytes, err := http.Post(
		host,
		http.Headers{},
		[]byte(url.QueryEscape(string(form.Encode()))),
		//[]byte("merchant_id%3Dtestmerchant123%26password%3Dtestpassword123%26customer_information%5Bfirst_name%5D%3DBarbara%26customer_information%5Bmiddle_name%5D%26customer_information%5Blast_name%5D%3DMiller%26customer_information%5Bemail%5D%3Dbarbara%40test.com%26customer_information%5Baddress1%5D%3D123%20Main%20St%26customer_information%5Baddress2%5D%26customer_information%5Bcity%5D%3DMiami%26customer_information%5Bprovince%5D%3DFL%26customer_information%5Bcountry%5D%3DUS%26customer_information%5Bpostal_code%5D%3D06614%26customer_information%5Bphone1%5D%3D88888888%26customer_information%5Bphone2%5D%3D55555555%26customer_information%5Bdob%5D%3D1983-06-02%26customer_information%5Bgender%5D%3DM%26customer_information%5Bid_values%5D%26user_name%3Dbarbara123%26user_number%3D123345%26reg_date%3D1982-06-01%26reg_ip_address%3D201.201.201.201%26reg_device_id%3D7ECS%2FdqbgMo4XWmeRtg2yaOufTP2X71GQsKJ1nMAoZbR%2Fq0uPQYKTYQVC3EtSz9N5lRmdBAbMv739lGFCme4qNd4J%2BIQsCgAq8vx8fjCInSwGD3OUDUFnT8tS2bDt%2Bo4qihK08ML6SX1SlYkG7YJ3wDZM2IR64zaPK61R9qc46cHP5%2B39QYMi5KQswRknT4eyY3sYT4UbVcjpPoWwQ4gaJAdIdMXRIZBUq2GRmr03CUNer%2F65MkncURpXwVmVS9Bk1gwcLQzgiHh6lAXMLzPUiRJ3KDtHsI9EkB8jJZaQ4eAFaXcB9q6fkoXdYJhlVol4nNw6ZsWjOGWu0leBcSuBqgdg3ffLo%2FG23lPrWKGSVovoGd9aSIT52Mu5Ps5hZgIjdxL535oouoXaAySJFMA%2FSihzKGUwR9jncrZBpshv%2BpPN%2FceQNXEhXp%2BoaA0WTDpsvJaZSehZYylIa1ScvaFdIWXbbqU%2BBlKi06heuQW3DQivjhuCR1TWO8o3jkjGBaeJcRdYF8kWwzRJnmH2ZqB%2Fn1gqr1EHpO0du98uLdptAC%2FPxkEOBF%2Baa8M%2FVJKJ75Rq%2F5e%2FncwDh4%2BZvvb8NzTXG9RtQFLR6cO0fzbGJLMGLszPLi%2BztpsDmHWpGgSYBvqK%2BaX%2B01jW67su%2FTHZfwVV68MeCBz0UoiJXsW9FDae63eR7f9rpkXbpBc6WY%2FF0trv51nlx0kDH6ZfbhoW9UwUlfO7DltZiP4j%2BcbGBrZpJO0xEDygGgqJjOdSNjd0X2ut3vmtAzAbdk88Bi7vHqWdkGRPGTIPbXm2TmEpYj9DOPkGAjbW0ABXVnmJRWFy36SzMUEaz2ZzIepr9o9hEvA8DQ%2FftFURjGJqUK62zeFy7m5nGY8O%2FDyPKCsjjI8fdPULV3Hi5CJ1UuojEAPPuvryeN1uIvV%2FCEZUzVUvmuPhez20JclP7DcmHmgJ%2Bt1z2nR4FFco55zWGi9Mml0nKA986q445vVEZgiOzRVjtVgyVV13g4X4axpAN4%2BkEqP78Gnkm3XdZzIqI8XiEzhvE%2BTbXXD65utTpBVOalT6lIzwXE1u1Jk%2BI6j%2FwPqe2f%2FW8gcEoMSHhYbz%2Bo9kxO%2BrN0vfXVbDgVarchzYvAYrkX6T%2F%2BaHcsFGPfoxqDH%2FxLrIjZqXEseRiVsuYFHEI00CA%2Bs68b414YG8jvPaLYo8JkHxdiK4M2ZEAVcRKIkKPHKoVB8PS8DSY1gz84jbq9w3AYCmuv2t5KE5eecsBOoFg95gehDt95tWAo30G%2BlT%2B8MJOjZZX46bTwiltQnEqVi4ehOUpIZuGPkpFi51SNOHpDVfTJ%2FPPO60sQKpCuefhGcPMu34XlQJ3aqe3lbr5tRR7WCb14F3CTRTMWuGA5ENs%2BKKuLOZiqXtF%2BZPIftoYUItJX%2Bx%2FCufrJFcyJb9zw3br2l3bNEq8qlU4tIal68iR%2B9LKwq9o%2BSRY%2B5CEnqC0BuxMSpBqSa72VfEgbKFRlLAId9GPj81sM9ueUd4mmC8Cb981lDlZGxE17NyUst1OxWHL4GWnVC5mQfpZZy%2BlxA1xIhCjjZFBQTmRtU8b%2FnNvo4%2BGNTf4zdshtK3frUma%2FxojNp8kr8wuVE2YCGMjQiy4EuD4M8xcTOLMno5yG8A3pEBsvAkeE6WMkEG58XWgOQ%2FNWKjV%2BSLIsaTkK6euRWvNe71Q069%2BdsM%2B9TW%2Blh9mmh4zuuTHQnEg7bJ1Dju4dOPB9ad1z3xCBEw4E8z42dVCFI%2BpTtgjk0WBO2kMHtML7PUfZf32%2BEGMzCTAZWtSTY7Ia%2F00X7ZrsP5rNSnCwPp0pXCEnG08wouMEWLfwNgmp6knJD68VowIaNmxZwjmY4AampNcOhTSwm54X2G%2FO7XQC8zjK23keVjEOe3oMjBPdwRbtJdZNzFrJeSWIt6iygNj5iCyVVnPlordjdIreYLZqkBLiilS0DJ441HGhzqWcNrp0SWI5yjoQJEF4zHcTmiz2AVpwO3JK7Cj9KrsAWNHoEhoORfUl74chbKkGQxwJpzcjmaJcehjkrq15ChOq1q93ZbsDRzKzHDRwRvd7ARMtlxuR%2FHkUeb36P4Qb%2FBu4HfISAHsxwTw8DOOrj83Nq8kCTPsaOe8C7DBU0C%2FidpwwIw77MpGzMOn5lTumIB6BH8VUMTavq8UJW4BzJIjYD73SpX5piEc6mOxBTp8jH3PY8hoESWDId9rMQcFPp%2FBdsHiKfohPKSmxOEtc6Jp212XiHi5EYoeLjQnW7b%2FsYkFXVW2LlV0X6OTjk3dT%2BJZQ7YXgMuDHDBvtQZVBRfToZgH4rPUAY18%2FOxDHo5hS0RLIWpVNkv2wbwmnh%2BQf0F6ZIlMU9zPNAEeVVw%2BNbq8xF94z2FeZcsFJJQgvkwCzyETzZMkmyHPXLQPXymwo3cQJhnoehM%2Fq8c%2FwvdhI0%2BqpIAH6IgLKEy%2BhemhA4o%2BYn5tgHDkDAh05Ei%2Bzncvq%2BbEIrEU5Jiay8V5J2%2BU%2FxttNOGt6adhOIorxlt8oajgSk5hE5ZTkYiwsPCeey00Rttw6MxrLUhfuO1dObWQnHQ46yE34in6g737P8jEnUYdo6X6TYxCH0d02GF343DQY%2FGl0xYSuomUDVeD0gmqPv4Q8nSgU7b3tZACVGfquf%2BU8Ns2J4iHA2gOcnKBo8AyW65afS2BhzOzzdiXHScBjOK%2FZfG4acz3uHZ0XqGZVQRrpaFtWAYY8k%2F87a6sHBLGFZVScxTUl0Tgbh46Nhhz0VHk2AsR1KyxqBtXril1kB5VAV2SNEuVIhhe%2Fws9L4G5akL%2F0PtoaAcXq4M7gFQgRIXxNh5YHUtOEz3q6QoJe6ieV8lMZe%2BbIBSCgMDxd0MrZiGMz0P9SY821JCshLjxD8%2Bj88nqf5HQo%2BV8sca%2Bf3asHe%2Fu0OzdfAlr1JhTiumk4JXkadDzYpApjgW4wNJl7cliiUB3EhN%2FF8sUwIeBUCbmjBYWZgm4fmfqXP2yXTpS%2FCBxnPhD9lbQO3ugqvi7XoUz4GNvzWOgAFTeuOm%2Bd97Naag7xmyCvjVBfKk27%2Fn5ryIG4yO0yrVzqEseNPwXNLxt4UbRaIqnqxf7fNE%2FgRM5XI3kvYOPGFS7OxJ3CuKYmfQC%2Flje4ws598iN2kVY6YeIGPJgP5zJX7cqT77xoWOCgid8KYlVqhA0DPExzUbcKyXgtoRHjCxPr2hvowAd1gFFQJQ%2BdF6YMmNvJqWlLRhb%2BU0Ti0Pin4S2kvd6QmCUCz2J%2BnnNi%2FzUSciuz713CfHvhDL1vko4ejKkG3%2BwNVqKpmLN8H6f%2BgDgu5qIFOnkqZuG0rKGJNvthbEaLetRged08nQ0Anh5IUIAzkk1bWZaR94PbMkKjWnkcxjvknbFPDEFVEXWvqqwMOpVnnNqTXF%2BsraFXTSTrPpazpMqCXBPl41c9XEVCx7NrtYhyQYZH%2F5rkAb6y2L0nTdPzaezWaoVa7m8iCUEOL9AA8JTJf3H%2FZ%2BtcI74NTOuc9xZbgx43e1DqCos5notI2TSvE0S55BIc9jEKimRthYbYQdVKmtUc9NOkgeExYOGZwMYpgl8sxuBWFYQGg7A8F5phO10SIEyiBaTTjsk0ZI31A3Mgqbmsv76Gy7oq%2FH8xREbysksh4ZE97%2F2CGaFYXJNxayWQWHWIfJGg%2F4msm%2FaD%2BiMCQIPFA5bvXVoBp%2BUlKv%2F7BkXORBZmKvz9vW4%2Ff3jSrA0rMmZof7UoUZ94iy5heh8r4dskEfz8BX8FZrtTiq8WDsMmW2RFuAAvvP5Sz7BJWOzyCLNSumWySAbBaQyU6dQ0vFV7LT4BqSve4I1f5AgeBAaG0Ndsk%2F1W7jSIGKTOUBd0Bg%2F0pCs5k2gW4crRJhyeuslloI6uq3Hbt3B7C4m9r8mRfBSqUH8pXyG6b1kuVfLbCQvVP7JtX8z0uKdS%2Bh1RcLAEoXlXmLByHsyW73q0H0r1te8uZWLctj0I6SVXFEFLQibuZXy0%2FgyXEispgsxT1BjQpf%2B1Zjzv59LDEHD%2Bqc%2FLWASZCtw3LKa68vAnHAQV01dPsqtGgoVMD1WjXGIV%2BsTsiJwEnfzYyG2Yu%2FxFSA%2FyD3SOUwJxbTNoY3hAYDx5joZrotn7eNHjoeKGcs0h%2BBcmgMQaNDGWwUoVrZXDiQWKswQcAH2Z2fX4gCHlpoNDmoOhdfEeG7EXMO5oxHM9hNqmOLDr%2Fw5jcWqpW0R%2BVWfNpumiw%2Bbe62slbjiKSR%2FtW5O7lBBzEUANduF%2BfwSE4t1M7P598V3Zi95No%2B%2Bpa9WdD%2FO3IwacAP%2BAPpfY3PdKDbCnslJk%3D%7C%7C%7CIII%7C%7C%7C2UjED88M5YNpGxqWSUJFtw%3D%3D%26device_fingerprint%3D01%3A01%3A01%3A01%3A01%3A01%3A01%26device_fingerprint_type%3D1%26bonus_code%3D25536%26bonus_amount%26bonus_submission_date%26site_skin_name%3Dwww.google.com%26how_did_you_hear%26affiliate_id%26pfc_status%3D2%26pfc_type%3D1"),
	)
	if err != nil {
		return nil, err
	}
	if responseStatus != 200 {
		err, _ := MapResponseError(responseStatus, responseBytes)
		return nil, err
	}

	response = &Response{}
	if err := json.Unmarshal(responseBytes, response); err != nil {

		// API returns invalid json:
		//{
		//	...
		//	"facts":[{
		//		"type":"1",
		//		"text":"Which one of the following addresses is associated with you?",
		//		"answers":[
		//			{
		//				"correct":"false",
		//				"text":"509 BIRDIE RD"
		//			},
		//			{
		//				"correct":"false",
		//				"text":"667 ASHWOOD NORT CT"
		//			},
		//			{
		//				"correct":"true",
		//				"text":"291 LYNCH RD"
		//			}
		//		} <--- "}" instead of "]"
		//	...
		//}
		// So let's try to extract valid part
		re := `(?si),\s+?"facts":.*$`
		match, _ := regexp.MatchString(re, string(responseBytes))
		if !match {
			return nil, err
		}

		if err := json.Unmarshal([]byte(regexp.MustCompile(re).ReplaceAllString(string(responseBytes), `}`)), response); err != nil {
			return nil, err
		}
	}

	if response.Status < 0 {
		response.Details = MapErrorCode(response.Status)
	}

	return response, nil
}

func (service service) AttachAuthPartToRequest(request *RegistrationRequest) *RegistrationRequest {
	request.MerchantID = service.config.MerchantID
	request.Password = service.config.Password

	return request
}

func generateUserNumber() (userNumber uuid.UUID, err error) {
	userNumber, err = uuid.NewV4()
	if err != nil {
		return uuid.UUID{}, err
	}

	return userNumber, nil
}

func convertModelToForm(model interface{}, result map[string]string){
	s := reflect.ValueOf(model)
	typeOfT := s.Type()

	if typeOfT == reflect.TypeOf(CustomerInformationField{}) {
		fName := s.Field(0).String()
		fVal := s.Field(1).String()

		if (len(fName) > 0) {
			result[fName] = fVal
		}
	} else {
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			convertModelToForm(f.Interface(), result)
		}
	}
}