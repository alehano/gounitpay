package gounitpay

import (
	"net/url"
	"testing"
)

var unitpay = New(Parameters{
	PublicKey:  "test",
	PrivateKey: "test",
})

func TestUnitpayNotif(t *testing.T) {
	paramStr := "method=pay&params[3ds]=1&params[account]=XXX&params[date]=2023-03-20 22:29:42&params[ip]=5.152.78.14&params[isPreauth]=0&params[operator]=card_not_rf&params[orderCurrency]=RUB&params[orderSum]=100.00&params[payerCurrency]=RUB&params[payerSum]=100.00&params[paymentType]=card&params[profit]=91.00&params[projectId]=XXX&params[purse]=XXX&params[signature]=XXXX&params[sum]=100&params[test]=0&params[unitpayId]=XXX"

	values, err := url.ParseQuery(paramStr)
	if err != nil {
		t.Errorf("bad query: %s", err)
	}

	notif, err := unitpay.ParseNotification(values)
	if err != nil {
		t.Error(err)
	}

	if notif.IsValidSignature() != true {
		t.Error("bad signature")
	}

}
