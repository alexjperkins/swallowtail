package handler

import (
	paymentsproto "swallowtail/s.payments/proto"
)

// PaymentsService ...
type PaymentsService struct {
	*paymentsproto.UnimplementedPaymentsServer
}
