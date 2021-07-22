package gounitpay

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type PaymentParameters struct {
	Account     string
	Description string
	Value       uint32

	Currency *string
	Locale   *string
	BackURL  *string
}

type Payment struct {
	*Unitpay
	*PaymentParameters
}

func (p *Payment) Signature() string {
	arguments := make([]string, 0)

	arguments = append(arguments, p.Account)
	if p.Currency != nil {
		arguments = append(arguments, *p.Currency)
	}
	arguments = append(arguments,
		p.Description,
		fmt.Sprintf("%d", p.Value),
		p.Unitpay.Parameters.PrivateKey,
	)

	return fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte(strings.Join(arguments, "{up}"))),
	)
}

func (p *Payment) QueryURL() string {
	return fmt.Sprintf(
		"%s?sum=%d&account=%s&desc=%s&signature=%s",
		p.Unitpay.Parameters.PublicKey,
		p.Value,
		p.Account,
		p.Description,
		p.Signature(),
	)
}
