package marshaling

import (
	"swallowtail/s.ftx/client"
	ftxproto "swallowtail/s.ftx/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func DepositDTOToProto(deposit *client.DepositRecord) *ftxproto.DepositRecord {
	return &ftxproto.DepositRecord{
		Coin:          deposit.Coin,
		Confirmations: deposit.Confirmations,
		ConfirmedTime: timestamppb.New(deposit.ConfirmedTime),
		Fee:           deposit.Fee,
		Id:            deposit.ID,
		SentTime:      timestamppb.New(deposit.SentTime),
		Size:          float32(deposit.Size),
		Status:        deposit.Status,
		Time:          timestamppb.New(deposit.Time),
		TransactionId: deposit.TXID,
	}
}

func DepositsDTOToProto(deposits []*client.DepositRecord) []*ftxproto.DepositRecord {
	protos := []*ftxproto.DepositRecord{}
	for _, d := range deposits {
		protos = append(protos, DepositDTOToProto(d))
	}

	return protos
}
