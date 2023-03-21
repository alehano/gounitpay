package gounitpay

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrNoRequiredArguments = errors.New("not all required arguments were provided")

	requiredNotificationParams = []string{
		"unitpayId", "projectId", "account", "orderSum", "orderCurrency",
	}
)

type Parameters struct {
	PublicKey  string
	PrivateKey string
}

type Unitpay struct {
	Parameters
}

func New(parameters Parameters) *Unitpay {
	return &Unitpay{
		Parameters: parameters,
	}
}

func (un *Unitpay) NewPayment(parameters PaymentParameters) *Payment {
	return &Payment{
		Unitpay:           un,
		PaymentParameters: &parameters,
	}
}

func (un *Unitpay) ParseNotification(formParameters url.Values) (*Notification, error) {
	var keys []string
	for k := range formParameters {
		if len(k) >= 9 && strings.Contains(k, "params") {
			keys = append(keys, k[7:len(k)-1])
		}
	}

	params := make(map[string]string)
	for _, k := range keys {
		params[k] = formParameters[fmt.Sprintf("params[%s]", k)][0]
	}

	for _, v := range requiredNotificationParams {
		_, ok := params[v]
		if !ok {
			return nil, ErrNoRequiredArguments
		}
	}

	notification := &Notification{
		Unitpay:        un,
		ID:             params["id"],
		ProjectID:      params["projectId"],
		Method:         formParameters.Get("method"),
		Type:           params["paymentType"],
		Account:        params["account"],
		PayerValue:     params["payerSum"],
		PayerCurrency:  params["payerCurrency"],
		OrderValue:     params["orderSum"],
		OrderCurrency:  params["orderCurrency"],
		Profit:         params["profit"],
		Phone:          params["phone"],
		Operator:       params["operator"],
		ThreeDS:        params["3ds"] == "1",
		SubscriptionID: params["subscriptionID"],
		Test:           params["test"] == "1",
		ErrorMessage:   params["errorMessage"],
		Date:           params["date"],
		Signature:      params["signature"],
	}

	delete(params, "signature")
	notification.params = params

	return notification, nil
}
