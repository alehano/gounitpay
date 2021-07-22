package gounitpay

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

type Notification struct {
	*Unitpay

	ID        uint32 // unitpayId
	ProjectID uint32 // projectId
	Method    string
	Type      string // paymentType

	Account       string
	PayerValue    uint32
	PayerCurrency string
	OrderValue    uint32
	OrderCurrency string
	Profit        uint32

	Phone    string
	Operator string

	ThreeDS        bool // 3ds
	SubscriptionID uint32
	Test           bool
	ErrorMessage   string
	Date           string

	params map[string]string
}

func (n *Notification) Signature() string {
	keys := make([]string, 0, len(n.params))
	for k := range n.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	arguments := make([]string, 0, len(n.params)+2)
	arguments = append(arguments, n.Method)
	for _, k := range keys {
		arguments = append(arguments, n.params[k])
	}
	arguments = append(arguments, n.Unitpay.PrivateKey)

	hash := sha256.New()
	return fmt.Sprintf(
		"%x",
		hash.Sum([]byte(strings.Join(arguments, "{up}"))),
	)
}
