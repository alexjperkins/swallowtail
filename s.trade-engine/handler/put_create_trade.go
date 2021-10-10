package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/s.trade-engine/marshaling"
)

// CreateTrade ...
func (s *TradeEngineService) CreateTrade(
	ctx context.Context, in *tradeengineproto.CreateTradeRequest,
) (*tradeengineproto.CreateTradeResponse, error) {
	// Validate trade.
	if err := validateTrade(in.Trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade", nil)
	}

	trade := marshaling.TradeProtoToDomain(in.Trade)
	errParams := map[string]string{
		"idempotency_key": trade.IdempotencyKey,
		"actor_id":        trade.ActorID,
	}

	// Idempotency check.
	alreadyExists, err := dao.TradeExists(ctx, trade.IdempotencyKey)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.failed_to_check_if_already_exists", errParams)
	}
	if alreadyExists {
		return nil, gerrors.AlreadyExists("failed_to_create_trade.already_exists", errParams)
	}

	// Create trade.
	if err := dao.CreateTrade(ctx, trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.dao", errParams)
	}

	// Read trade back out; we don't know the internal uuid, so we use our idempotency key which
	// guranteed to be unique.
	embelishedTrade, err := dao.ReadTradeByIdempotencyKey(ctx, trade.IdempotencyKey)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.failed_to_read_created_trade_back", errParams)
	}

	return &tradeengineproto.CreateTradeResponse{
		TradeId: embelishedTrade.TradeID,
		Created: timestamppb.New(embelishedTrade.Created),
	}, nil
}

func validateTrade(trade *tradeengineproto.Trade) error {
	switch {
	case trade == nil:
		return gerrors.BadParam("missing_param.trade", nil)
	case trade.Asset == "":
		return gerrors.BadParam("missing_param.asset", nil)
	case trade.IdempotencyKey == "":
		return gerrors.BadParam("missing_param.idempotency_key", nil)
	case trade.ActorId == "":
		return gerrors.BadParam("missing_param.actor_id", nil)
	case len(trade.Entries) == 0:
		return gerrors.BadParam("missing_param.entries", nil)
	case trade.TradeType == tradeengineproto.TRADE_TYPE_FUTURES_PERPETUALS && trade.StopLoss == 0:
		return gerrors.FailedPrecondition("missing_param.stoploss_cannot_be_zero_for_futures_perpetuals", nil)
	}

	for _, entry := range trade.Entries {
		if entry == 0 {
			return gerrors.BadParam("bad_param.zero_valued_entry", nil)
		}
	}

	return nil
}
