package gounitpay

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

type Notification struct {
	*Unitpay

	ID        string // unitpayId
	ProjectID string // projectId
	Method    string
	Type      string // paymentType

	Account       string
	PayerValue    string
	PayerCurrency string
	OrderValue    string
	OrderCurrency string
	Profit        string

	Phone    string
	Operator string

	ThreeDS        bool // 3ds
	SubscriptionID string
	Test           bool
	ErrorMessage   string
	Date           string
	Signature      string

	params map[string]string
}

func (n *Notification) MakeSignature() string {
	keys := make([]string, 0, len(n.params))
	for k := range n.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	arguments := make([]string, 0, len(n.params)+2)
	arguments = append(arguments, n.Method)
	for _, k := range keys {
		if n.params[k] == "" {
			continue
		}
		arguments = append(arguments, n.params[k])
	}
	arguments = append(arguments, n.Unitpay.PrivateKey)

	argStr := strings.Join(arguments, "{up}")

	hash := sha256.New()
	hash.Write([]byte(argStr))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (n *Notification) IsValidSignature() bool {
	return n.Signature == n.MakeSignature()
}
