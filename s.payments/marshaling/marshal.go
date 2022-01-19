package marshaling

import (
	"swallowtail/s.payments/domain"
	paymentsproto "swallowtail/s.payments/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentsToProtos ...
func PaymentsToProtos(in []*domain.Payment) []*paymentsproto.Payment {
	var pp = make([]*paymentsproto.Payment, 0, len(in))
	for _, p := range in {
		pp = append(pp, &paymentsproto.Payment{
			PaymentTimestamp: timestamppb.New(p.Timestamp),
			TransactionId:    p.TransactionID,
			AmountInUsdt:     float32(p.AmountInUSDT),
		})
	}

	return pp
}
