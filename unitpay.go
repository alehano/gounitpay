package gounitpay

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrNoRequiredArguments = errors.New("not all required arguments were provided")
	ErrBadArgumentFormat   = errors.New("bad argument format")
	ErrBadSignature        = errors.New("bad signature")

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

func (un *Unitpay) ParseNotification(method string, formParameters url.Values) (*Notification, error) {
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

	signature := params["signature"]

	// Не участвуют в формировании подписи.
	delete(params, "sign")
	delete(params, "signature")

	id, err := strconv.ParseUint(params["unitpayId"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: unitpayId", ErrBadArgumentFormat)
	}

	projectId, err := strconv.ParseUint(params["projectId"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: projectId", ErrBadArgumentFormat)
	}

	payerSum, err := strconv.ParseUint(params["payerSum"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: payerSum", ErrBadArgumentFormat)
	}

	orderSum, err := strconv.ParseUint(params["orderSum"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: orderSum", ErrBadArgumentFormat)
	}

	profit, err := strconv.ParseUint(params["profit"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: profit", ErrBadArgumentFormat)
	}

	subscriptionID, err := strconv.ParseUint(params["subscriptionId"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: subscriptionId", ErrBadArgumentFormat)
	}

	notification := &Notification{
		Unitpay:        un,
		ID:             uint32(id),
		ProjectID:      uint32(projectId),
		Method:         params["method"],
		Type:           params["paymentType"],
		Account:        params["account"],
		PayerValue:     uint32(payerSum),
		PayerCurrency:  params["payerCurrency"],
		OrderValue:     uint32(orderSum),
		OrderCurrency:  params["orderCurrency"],
		Profit:         uint32(profit),
		Phone:          params["phone"],
		Operator:       params["operator"],
		ThreeDS:        params["3ds"] == "1",
		SubscriptionID: uint32(subscriptionID),
		Test:           params["test"] == "1",
		ErrorMessage:   params["errorMessage"],
		Date:           params["date"],
	}

	if signature != notification.Signature() {
		return nil, ErrBadSignature
	}

	return notification, nil
}
