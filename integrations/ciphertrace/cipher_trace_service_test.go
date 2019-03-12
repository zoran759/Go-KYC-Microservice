package ciphertrace

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var testCipherService *CipherService
var testWalletStruct *Wallet
var transactionHashes []string

func TestMain(m *testing.M) {
	testCipherService = NewCipherService("https://rest.ciphertrace.com", "a14b5221f82ffef1b7c75ef5a44a7e0186fd0e90d2b2eb0830307eccaeaf02b9", "mpinderi_key1")
	os.Exit(m.Run())
}

func TestCipherService_GetWalletByAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "success", address: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt"},
		{name: "error", address: "testerror"},
		{name: "empty", address: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GetWalletByAddress(test.address)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success address. Err: %v", err)
				}
				if res.WalletID == "" {
					t.Error("Empty Wallet Id")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				t.Log(res.WalletID)
				testWalletStruct = res
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GetWalletByWalletId(t *testing.T) {
	tests := []struct {
		name     string
		walletId string
	}{
		{name: "success", walletId: "094f8d86"},
		{name: "error", walletId: "testerror"},
		{name: "empty", walletId: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GetWalletByWalletId(test.walletId)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success wallet Id. Err: %v", err)
				}
				if res.WalletID == "" {
					t.Error("Empty Wallet Id")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GetWalletWithAddresses(t *testing.T) {
	tests := []struct {
		name     string
		walletId string
		count    int
		offset   int
	}{
		{name: "success", walletId: "094f8d86", count: 100, offset: 100},
		{name: "error", walletId: "testerror", count: 100, offset: 100},
		{name: "empty", walletId: "", count: 100, offset: 100},
		{name: "zeroCount", walletId: "094f8d86", count: 0, offset: 100},
		{name: "biggerCount", walletId: "094f8d86", count: 100000, offset: 100},
		{name: "notMultipliers", walletId: "094f8d86", count: 20, offset: 210},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GetWalletWithAddresses(test.walletId, test.offset, test.count)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success wallet Id. Err: %v", err)
				}
				if res.WalletID == "" {
					t.Error("Empty Wallet Id")
				}
				if len(res.Addresses) != test.count {
					fmt.Println(len(res.Addresses))
					t.Errorf("Addressess returned less then count")
				}
				if res.AddressOffset != test.offset {
					fmt.Println(res.AddressOffset)
					t.Errorf("Offset not equal to sent offset")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GetTransactionHistoryForAddress(t *testing.T) {
	tests := []struct {
		name      string
		addressId string
		startDate int64
		endDate   int64
	}{
		{name: "success", addressId: "3GVbgAcs7Gr6q118i8FYYVKnshsQQwwpUK", startDate: time.Now().AddDate(-3, 0, 0).Unix(), endDate: time.Now().Unix()},
		{name: "empty", addressId: "", startDate: time.Now().AddDate(0, -1, 0).Unix(), endDate: time.Now().Unix()},
		{name: "moreOpts", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(-1, 0, 0).Unix(), endDate: time.Now().Unix()},
		{name: "biggerStartDate", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(1, 0, 0).Unix(), endDate: time.Now().Unix()},
		{name: "success", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(-3, 0, 0).Unix()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &TransactionHistoryForAddress{}
			var err error
			var sts int
			if test.name == "moreOpts" {
				res, err, sts = testCipherService.GetTransactionHistoryForAddress(test.addressId, test.startDate, test.endDate, time.Now().Unix())
			} else {
				res, err, sts = testCipherService.GetTransactionHistoryForAddress(test.addressId, test.startDate, test.endDate)
			}
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success wallet Id. Err: %v", err)
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				if res.Address == "" {
					t.Error("Empty Wallet Id")
				}
				if len(res.Transactions) < 1 {
					t.Errorf("Transactions did not return")
				}
				transactionHashes = res.Transactions
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GetTransactionsHistoryByTxHash(t *testing.T) {
	tests := []struct {
		name     string
		txhashes []string
	}{
		//{name: "success", txhashes:transactionHashes[:2]},
		{name: "success", txhashes: []string{"a9f6fee0137b021e9637fd9faaeec740533cd64c7720444840ad19e883f40b27"}},
		{name: "errorToMany", txhashes: []string{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "}},
		{name: "errorNone"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GetTransactionsHistoryByTxHash(test.txhashes...)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success wallet Id. Err: %v", err)
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				if len(res.Transactions) < 1 {
					t.Errorf("no transactions returned")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_SearchAddressInfo(t *testing.T) {
	tests := []struct {
		name      string
		addressId string
		features  []string
		startDate int64
		endDate   int64
	}{
		{name: "success", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(-5, 0, 0).Unix(), endDate: time.Now().Unix(), features: []string{"tx", "ip"}},
		{name: "empty", addressId: "", startDate: time.Now().AddDate(0, -1, 0).Unix(), endDate: time.Now().Unix()},
		{name: "moreOpts", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(-1, 0, 0).Unix(), endDate: time.Now().Unix()},
		{name: "biggerStartDate", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(1, 0, 0).Unix(), endDate: time.Now().Unix()},
		{name: "success", addressId: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt", startDate: time.Now().AddDate(-5, 0, 0).Unix(), features: []string{"tx"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &AddressSearchesInfo{}
			var err error
			var sts int
			if test.name == "moreOpts" {
				res, err, sts = testCipherService.SearchAddressInfo(test.addressId, test.features, test.startDate, test.endDate, time.Now().Unix())
			} else {
				res, err, sts = testCipherService.SearchAddressInfo(test.addressId, test.features, test.startDate, test.endDate)
			}
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success wallet Id. Err: %v", err)
				}
				if res.Address == "" {
					t.Error("Empty Wallet Id")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				if len(res.TxHistory) < 1 {
					t.Errorf("Length of transactions is 0")
				}
				t.Log(res)
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GtAddressRiskInfo(t *testing.T) {
	tests := []struct {
		name   string
		txhash string
	}{
		//{name: "success", txhash:transactionHashes[1]},
		{name: "success", txhash: "a9f6fee0137b021e9637fd9faaeec740533cd64c7720444840ad19e883f40b27"},
		{name: "errorNone", txhash: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GtAddressRiskInfo(test.txhash)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success address Id. Err: %v", err)
				}
				if len(res.AddressRisks) < 1 {
					t.Errorf("no transactions returned")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GtAddressRiskInfoETH(t *testing.T) {
	tests := []struct {
		name   string
		txhash string
	}{
		{name: "success", txhash: "0xc9d34946894f939770d87e2e6fd39825c13b5d9a08c20ca73b983304e080f9d5"},
		{name: "errorNone", txhash: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GtAddressRiskInfoETH(test.txhash)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success address Id. Err: %v", err)
				}
				if len(res.AddressRisks) < 1 {
					t.Errorf("no transactions returned")
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_SingleAddressRiskInfo(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "success", address: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt"},
		{name: "errorNone", address: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GtSingleAddressRiskInfo(test.address)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success address Id. Err: %v", err)
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				if res.Address == "" {
					t.Errorf("Address empty.")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}

func TestCipherService_GtSingleAddressRiskInfoETH(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "success", address: "39JBmVmirpsnxJG7ZM38YvjdSghkpQHURt"},
		{name: "errorNone", address: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err, sts := testCipherService.GtSingleAddressRiskInfo(test.address)
			if test.name == "success" {
				if err != nil {
					t.Errorf("Error with success address Id. Err: %v", err)
				}
				if sts != 200 {
					t.Errorf("Status should be 200")
				}
				if res.Address == "" {
					t.Errorf("Address empty.")
				}
			} else {
				if err == nil {
					t.Error("Test passed with wrong data")
				}
			}
		})
	}
}
