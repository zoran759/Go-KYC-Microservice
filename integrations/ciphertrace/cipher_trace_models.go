package ciphertrace

import "github.com/shopspring/decimal"

type Wallet struct {
	WalletID string `json:"walletId"`
	Owner    struct {
		Name        string `json:"name"`
		Subpoenable bool   `json:"subpoenable"`
		URL         string `json:"url"`
		Country     string `json:"country"`
		Type        string `json:"type"`
	} `json:"owner"`
	TotalAddressCount int `json:"totalAddressCount"`
	Revision          int `json:"revision"`
}

type WalletWithAddresses struct {
	AddressOffset int `json:"addressOffset"`
	Wallet
	Addresses []string `json:"addresses"`
}

type TransactionHistoryForAddress struct {
	Address      string   `json:"address"`
	StartDate    int      `json:"startDate"`
	EndDate      int      `json:"endDate"`
	Transactions []string `json:"transactions"`
}

type Transaction struct {
	TxHash  string `json:"tx_hash"`
	Outputs []struct {
		Pos     int     `json:"pos"`
		Address string  `json:"address"`
		Value   float64 `json:"value"`
	} `json:"outputs"`
	Total  float64 `json:"total"`
	Inputs []struct {
		Pos     int     `json:"pos"`
		Address string  `json:"address"`
		Value   float64 `json:"value"`
	} `json:"inputs"`
	Date int     `json:"date"`
	Fee  float64 `json:"fee"`
}

type IpHistory struct {
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Country       string  `json:"country"`
	Longitude     float64 `json:"longitude"`
	Date          int     `json:"date"`
	IPAddress     string  `json:"ipAddress"`
	ClientVersion string  `json:"clientVersion"`
}

type TransactionDetails struct {
	Transactions []Transaction          `json:"transactions"`
	Addresses    map[string]Wallet      `json:"addresses"`
	IpHistory    map[string][]IpHistory `json:"ipHistory"`
}

type AddressSearchesInfo struct {
	LastUsedBlockHeight int         `json:"lastUsedBlockHeight"`
	QuerySpent          float64     `json:"querySpent"`
	QueryEndingBalance  float64     `json:"queryEndingBalance"`
	EndDate             int         `json:"endDate"`
	TotalSpendCount     int         `json:"totalSpendCount"`
	TotalSpent          float64     `json:"totalSpent"`
	TotalDepositCount   int         `json:"totalDepositCount"`
	QueryDeposits       float64     `json:"queryDeposits"`
	CurrentBalance      float64     `json:"currentBalance"`
	QueryDepositCount   int         `json:"queryDepositCount"`
	IPHistory           []IpHistory `json:"ipHistory"`
	QuerySpendCount     int         `json:"querySpendCount"`
	Address             string      `json:"address"`
	TxHistory           []struct {
		TxHash   string  `json:"txHash"`
		TxIndex  int     `json:"txIndex"`
		Balance  float64 `json:"balance"`
		Date     int     `json:"date"`
		Received float64 `json:"received"`
		Spent    float64 `json:"spent"`
	} `json:"txHistory"`
	InCase        bool    `json:"inCase"`
	StartDate     int     `json:"startDate"`
	Wallet        Wallet  `json:"wallet"`
	TotalDeposits float64 `json:"totalDeposits"`
}

type Risk struct {
	OutputValue     decimal.Decimal `json:"outputValue"`
	CallBackSeconds int             `json:"callBackSeconds"`
	Risk            decimal.Decimal `json:"risk"`
	Address         string          `json:"address"`
	InputValue      decimal.Decimal `json:"inputValue"`
}

type AddressRisk struct {
	CallBackSeconds int             `json:"callBackSeconds"`
	Risk            float64         `json:"risk"`
	Txhash          string          `json:"txhash"`
	AddressRisks    map[string]Risk `json:"addressRisks"`
	UpdatedToBlock  int             `json:"updatedToBlock"`
}

type SingleAddressRisk struct {
	Address         string  `json:"address"`
	Risk            float64 `json:"risk"`
	UpdatedToBlock  int     `json:"updatedToBlock"`
	CallBackSeconds int     `json:"callBackSeconds"`
}

type ETHAddressRisk struct {
	CallBackSeconds int     `json:"callBackSeconds"`
	Address         string  `json:"address"`
	Risk            float64 `json:"risk"`
	UpdateToBlock   int     `json:"updateToBlock"`
	Balance         string  `json:"balance"`
}
