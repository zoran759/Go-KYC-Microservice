package ciphertrace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CipherService struct {
	endpoint   string
	privateKey string
	username   string
}

//Create a new CipherTrace Service by providing the endpoint, private key and username
func NewCipherService(endpoint, privateKey, username string) *CipherService {
	return &CipherService{
		endpoint:   endpoint,
		privateKey: privateKey,
		username:   username,
	}
}

//Get a wallet object from the address id
func (cipherService *CipherService) GetWalletByAddress(address string) (*Wallet, error, int) {
	// if the provided address is empty return a 400 status and a message
	if address == "" {
		return &Wallet{}, errors.New("address can not be empty"), http.StatusBadRequest
	}
	client := http.Client{}                                                             //Initialize a new http client
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers
	url := fmt.Sprintf("%v/api/v1/wallet?address=%v", cipherService.endpoint, address)  // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload
	if err != nil {
		return &Wallet{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header
	res, err := client.Do(req)                         // Perform the request
	if err != nil {
		return &Wallet{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}

	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &Wallet{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}

	wallet := &Wallet{} // Create a new wallet object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(wallet) // Decode the response into the wallet object
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return wallet, err, http.StatusOK // Return the response and the 200 status
}

// Get a wallet object by wallet id
func (cipherService *CipherService) GetWalletByWalletId(walletId string) (*Wallet, error, int) {
	// if the provided walled id is empty return a 400 status and a message

	if walletId == "" {
		return &Wallet{}, errors.New("wallet id can not be empty"), http.StatusBadRequest
	}
	client := http.Client{}                                                               // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey)   // Create the authentication string to be put on headers
	url := fmt.Sprintf("%v/api/v1/wallet?wallet_id=%v", cipherService.endpoint, walletId) // Create the url string
	req, err := http.NewRequest("GET", url, nil)                                          // Create a new request by providing the method, url and payload
	if err != nil {
		return &Wallet{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header
	res, err := client.Do(req)                         // Perform the request
	if err != nil {
		return &Wallet{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}
	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &Wallet{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	wallet := &Wallet{} // Create a new wallet object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(wallet) // Return the response and the 200 status
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return wallet, err, http.StatusOK // Decode the response into the wallet object

}

// Get a wallet object with all its addresses by specifying walletId, offset and count. Offset and count should be multipliers of 100 and the count should not be 0
func (cipherService *CipherService) GetWalletWithAddresses(walletId string, offset int, count int) (*WalletWithAddresses, error, int) {
	if walletId == "" {
		return &WalletWithAddresses{}, errors.New("wallet id can not be empty"), http.StatusBadRequest // if the provided walled id is empty return a 400 status and a message
	}
	if count == 0 {
		return &WalletWithAddresses{}, errors.New("count should be bigger than 0"), http.StatusBadRequest // if the provided count is 0 return a 400 status and a message
	}
	if offset%100 != 0 || count%100 != 0 {
		return &WalletWithAddresses{}, errors.New("offset and counter should be multipliers of 100"), http.StatusBadRequest // if the provided count or offest  is not a multiplier of 100 return a 400 status and a message
	}
	if count > 10000 {
		return &WalletWithAddresses{}, errors.New("count should not be bigger than 10000"), http.StatusBadRequest // if the provided count is bigger than 10000 return a 400 status and a message
	}
	client := http.Client{}                                                                                                           // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey)                                               // Create the authentication string to be put on headers
	url := fmt.Sprintf("%v/api/v1/wallet/addresses?wallet_id=%v&count=%v&offset=%v", cipherService.endpoint, walletId, count, offset) // Create the url string
	req, err := http.NewRequest("GET", url, nil)                                                                                      // Create a new request by providing the method, url and payload

	if err != nil {
		return &WalletWithAddresses{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header

	res, err := client.Do(req) // Perform the request
	if err != nil {
		return &WalletWithAddresses{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}

	// If the status code is not 200 extract the error message from the response and return it

	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &WalletWithAddresses{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}

	wallet := &WalletWithAddresses{} // Create a new wallet with address object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(wallet) // Decode the response into the wallet object
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return wallet, err, http.StatusOK // Return the response and the 200 status

}

// Get transaction history for an address by specifying the addressId and OPTIONALLY start date and end date
func (cipherService *CipherService) GetTransactionHistoryForAddress(addressId string, opts ...int64) (*TransactionHistoryForAddress, error, int) {
	if addressId == "" {
		return &TransactionHistoryForAddress{}, errors.New("address id can not be empty"), http.StatusBadRequest // if the provided address id is empty return a 400 status and a message
	}
	if len(opts) > 2 {
		return &TransactionHistoryForAddress{}, errors.New("send only start date and end date in options"), http.StatusBadRequest // if the provided length of options is bigger than 2 return a 400 status and a message
	}
	// The following statements are used to parse the date from the options
	var url string
	var startDate int64
	var endDate int64
	now := time.Now().Unix()
	withDate := true
	if len(opts) == 1 {
		startDate = opts[0]
		endDate = time.Now().Unix()
	} else if len(opts) == 2 && opts[0] != 0 {
		if opts[1] == 0 {
			startDate = opts[0]
			endDate = time.Now().Unix()
		} else {
			startDate = opts[0]
			endDate = opts[1]
		}
	} else {
		withDate = false
	}
	// Check if start date and end date exceed the current time and return a 400 status and a message
	if startDate > now || endDate > now {
		return &TransactionHistoryForAddress{}, errors.New("start date and end date should be before  now"), http.StatusBadRequest
	}
	client := http.Client{} // Initialize a new client instance

	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create a new auth string

	// Check weather the optional dates are provided and create the request url accordingly
	if withDate {
		fmt.Println(startDate, endDate)
		url = fmt.Sprintf("%v/api/v1/tx/search?address=%v&startdate=%v&enddate=%v", cipherService.endpoint, addressId, startDate, endDate)
	} else {
		url = fmt.Sprintf("%v/api/v1/tx/search?address=%v", cipherService.endpoint, addressId)
	}

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &TransactionHistoryForAddress{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message

	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header

	res, err := client.Do(req) // Perform the request

	if err != nil {
		return &TransactionHistoryForAddress{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}
	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &TransactionHistoryForAddress{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	transHistory := &TransactionHistoryForAddress{} // Create a new transaction history object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(transHistory) // Decode the response into the wallet object
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return transHistory, err, http.StatusOK // Return the response and the 200 status

}

// Get the transaction history by specifying an array of txhashes
func (cipherService *CipherService) GetTransactionsHistoryByTxHash(txhashes ...string) (*TransactionDetails, error, int) {
	if len(txhashes) < 1 {
		return &TransactionDetails{}, errors.New("you should specify at least 1 transaction id "), http.StatusBadRequest // If the length of tx hashes provided is less than one return error and 400 status
	}
	if len(txhashes) > 10 {
		return &TransactionDetails{}, errors.New("you should specify at most 10 transaction ids "), http.StatusBadRequest // If the length of tx hashes is bigger than 10 return 400 and an error
	}
	client := http.Client{}                                                             // Initialize a new client instance from package http
	txhashesstr := strings.Join(txhashes, ",")                                          // Join all the provided hashes into a string using ',' as a delimiter
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers
	url := fmt.Sprintf("%v/api/v1/tx?txhashes=%v", cipherService.endpoint, txhashesstr) // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &TransactionDetails{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header

	res, err := client.Do(req) // Perform the request
	if err != nil {
		return &TransactionDetails{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}

	// If the status code is not 200 extract the error message from the response and return it

	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &TransactionDetails{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	transDetails := &TransactionDetails{} // Create a new transaction details object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(transDetails) // Decode the response into the transaction details object

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return transDetails, err, http.StatusOK // Return the response and the 200 status

}

// This query returns all information regarding an Address. Current balance information as well as (optional if specified as features in the features array) balance history with transaction hashes and IP Address match history.
// You can also OPTIONALLY specify start date and end date.

func (cipherService *CipherService) SearchAddressInfo(addressId string, features []string, dates ...int64) (*AddressSearchesInfo, error, int) {
	if addressId == "" {
		return &AddressSearchesInfo{}, errors.New("address id can not be empty"), http.StatusBadRequest // If the provided address id is empty return a 400 status code and an error message
	}

	if len(dates) > 2 {
		return &AddressSearchesInfo{}, errors.New("send only start date and end date in options"), http.StatusBadRequest // If the length of the provided optional dates is more than 2 return an 400 status and an error message
	}
	featuresStr := strings.Join(features, ",") // Join the optional features by ',' as delimiter

	//The following statements extract the date from the optional variadic argument dates
	var url string
	var startDate int64
	var endDate int64
	now := time.Now().Unix()
	withDate := true
	if len(dates) == 1 {
		startDate = dates[0]
		endDate = time.Now().Unix()
	} else if len(dates) == 2 && dates[0] != 0 {
		if dates[1] == 0 {
			startDate = dates[0]
			endDate = time.Now().Unix()
		} else {
			startDate = dates[0]
			endDate = dates[1]
		}
	} else {
		withDate = false
	}
	// If the start date or the end date are latter than the current date a 400 status is returned and an error message
	if startDate > now || endDate > now {
		return &AddressSearchesInfo{}, errors.New("start date and end date should be before  now"), http.StatusBadRequest
	}
	client := http.Client{} // Create a new client instance

	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the auth url to be put on headers

	// Assign the url accordingly weather the optional dates are provided or not
	if withDate {
		url = fmt.Sprintf("%v/api/v1/address/search?address=%v&features=%v&startdate=%v&enddate=%v", cipherService.endpoint, addressId, featuresStr, startDate, endDate)
	} else {
		url = fmt.Sprintf("%v/api/v1/address/search?address=%v&features=%v", cipherService.endpoint, addressId, featuresStr)
	}

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &AddressSearchesInfo{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header
	req.Header.Add("Authorization", auth)              // Add the authorization header

	res, err := client.Do(req) // Perform the request
	if err != nil {
		return &AddressSearchesInfo{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message
	}

	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &AddressSearchesInfo{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}

	// Create a new address search info object
	transHistory := &AddressSearchesInfo{}

	decoder := json.NewDecoder(res.Body) //Create a new decoder

	err = decoder.Decode(transHistory) // Decode the response into the address search info object

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return transHistory, err, http.StatusOK // Return the response and the 200 status

}

// Get the address Risk info for bitcoin by specifying txhash
func (cipherService *CipherService) GtAddressRiskInfo(txhash string) (*AddressRisk, error, int) {
	if txhash == "" {
		return &AddressRisk{}, errors.New("transaction id can not be empty"), http.StatusBadRequest // if the tx hash is not provided return a 400 bad request and an error message
	}
	client := http.Client{}                                                             // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers

	url := fmt.Sprintf("%v/aml/v1/btc/risk?txhash=%v", cipherService.endpoint, txhash) // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &AddressRisk{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message
	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header

	req.Header.Add("Authorization", auth) // Add the authorization header

	res, err := client.Do(req) // Perform the request

	if err != nil {
		return &AddressRisk{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message

	}
	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &AddressRisk{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	riskDetails := &AddressRisk{} // Create a new address risk object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(riskDetails) // Decode the response

	if err != nil {
		return riskDetails, err, http.StatusInternalServerError
	}

	return riskDetails, nil, http.StatusOK // If there was no error return response and a 200 status

}

// Get a single address risk info by specifying the address id for Bitcoin
func (cipherService *CipherService) GtSingleAddressRiskInfo(addressId string) (*SingleAddressRisk, error, int) {
	if addressId == "" {
		return &SingleAddressRisk{}, errors.New("transaction id can not be empty"), http.StatusBadRequest // If the address id is empty return a 400 status and an error
	}
	client := http.Client{}                                                             // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers

	url := fmt.Sprintf("%v/aml/v1/btc/risk?address=%v", cipherService.endpoint, addressId) // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &SingleAddressRisk{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message

	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header

	req.Header.Add("Authorization", auth) // Add the authorization header

	res, err := client.Do(req) // Perform the request

	if err != nil {
		return &SingleAddressRisk{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message

	}
	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &SingleAddressRisk{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	riskDetails := &SingleAddressRisk{} // Create a new single risk address object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(riskDetails) // Decode the response into the single risk address object

	if err != nil {
		return riskDetails, err, http.StatusInternalServerError
	}
	return riskDetails, err, http.StatusOK // Return the response and the 200 status

}

// Get a single address risk info by specifying the address id for ETH

func (cipherService *CipherService) GtSingleAddressRiskInfoETH(addressId string) (*ETHAddressRisk, error, int) {
	if addressId == "" {
		return &ETHAddressRisk{}, errors.New("transaction id can not be empty"), http.StatusBadRequest // If the address id is empty return a 400 status and an error
	}
	client := http.Client{}                                                             // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers

	url := fmt.Sprintf("%v/aml/v1/eth/risk?address=%v", cipherService.endpoint, addressId) // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &ETHAddressRisk{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message

	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header

	req.Header.Add("Authorization", auth) // Add the authorization header

	res, err := client.Do(req) // Perform the request
	if err != nil {
		return &ETHAddressRisk{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message

	}
	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &ETHAddressRisk{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	riskDetails := &ETHAddressRisk{} // Create a new ETH ris info object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(riskDetails) // Decode the response into the wallet object

	if err != nil {
		return riskDetails, err, http.StatusInternalServerError
	}
	return riskDetails, err, http.StatusOK // Return the response and the 200 status

}

// Get the address Risk info for ETH by specifying txhash
func (cipherService *CipherService) GtAddressRiskInfoETH(txhash string) (*AddressRisk, error, int) {
	if txhash == "" {
		return &AddressRisk{}, errors.New("transaction id can not be empty"), http.StatusBadRequest // If the tx hash is empty return a 400 status and an error
	}
	client := http.Client{}                                                             // Create a new client instance
	auth := fmt.Sprintf("ctv1:%v:%v", cipherService.username, cipherService.privateKey) // Create the authentication string to be put on headers

	url := fmt.Sprintf("%v/aml/v1/eth/risk?txhash=%v", cipherService.endpoint, txhash) // Create the url string

	req, err := http.NewRequest("GET", url, nil) // Create a new request by providing the method, url and payload

	if err != nil {
		return &AddressRisk{}, err, http.StatusInternalServerError // If new request returns an error return 500 status and the error message

	}
	req.Header.Add("Content-Type", "application/json") // Add application/json header

	req.Header.Add("Authorization", auth) // Add the authorization header

	res, err := client.Do(req) // Perform the request

	if err != nil {
		return &AddressRisk{}, err, http.StatusInternalServerError // If the request returns an error return 500 status and the error message

	}

	// If the status code is not 200 extract the error message from the response and return it
	if res.StatusCode != 200 {
		var resByte []byte
		if res.Body != nil {
			resByte, _ = ioutil.ReadAll(res.Body)
		}
		return &AddressRisk{}, errors.New("status not 200. error: " + string(resByte)), res.StatusCode
	}
	riskDetails := &AddressRisk{} // Create a new address risk object

	decoder := json.NewDecoder(res.Body) // Create a new decoder

	err = decoder.Decode(riskDetails) // Decode the response into the address risk object

	if err != nil {
		return riskDetails, err, http.StatusInternalServerError
	}
	return riskDetails, nil, http.StatusOK // Return the response and the 200 status

}
