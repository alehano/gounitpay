package gounitpay

import (
	"testing"
)

var unitpay = New(Parameters{
	PublicKey:  "428369-e60f0",
	PrivateKey: "be5ab023368bbfa5a830c7701e14f9d3",
})

func TestUnitpay_NewPayment(t *testing.T) {
	payment := unitpay.NewPayment(PaymentParameters{
		Account:     "12345",
		Description: "desc",
		Value:       125,
		Currency:    ptrS("RUB"),
	})

	neededSignature := "8d0dfb5d248a9b38cf171ef99f6d48be0d34d0f3bde193653f948d901cb2b0b4"
	if payment.Signature() != neededSignature {
		t.Errorf(
			"bad signature (%s != %s)",
			payment.Signature(),
			neededSignature,
		)
	}

	payment = unitpay.NewPayment(PaymentParameters{
		Account:     "123452",
		Description: "desc2",
		Value:       126,
		Currency:    ptrS("RUB2"),
	})

	neededSignature = "3484e2d694c94ee3de2ebec06bd4f612243e151ab1588e3554838b1edfe55adb"
	if payment.Signature() != neededSignature {
		t.Errorf(
			"bad signature (%s != %s)",
			payment.Signature(),
			neededSignature,
		)
	}
}

func ptrS(s string) *string {
	return &s
}
