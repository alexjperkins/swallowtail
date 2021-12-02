package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.trade-engine/dao"
	tradeengineproto "swallowtail/s.trade-engine/proto"

	"google.golang.org/protobuf/types/known/timestamppb"

	"swallowtail/s.trade-engine/marshaling"
)

// CreateTradeStrategy creates a trade strategy for participants to execute.
func (s *TradeEngineService) CreateTradeStrategy(
	ctx context.Context, in *tradeengineproto.CreateTradeStrategyRequest,
) (*tradeengineproto.CreateTradeStrategyResponse, error) {
	// Validate trade strategy.
	if err := validateTradeStrategy(in.TradeStrategy); err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade", nil)
	}

	trade := marshaling.TradeStrategyProtoToDomain(in.TradeStrategy)
	errParams := map[string]string{
		"idempotency_key": trade.IdempotencyKey,
		"actor_id":        trade.ActorID,
	}

	// Idempotency check.
	alreadyExists, err := dao.TradeStrategyExists(ctx, trade.IdempotencyKey)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.failed_to_check_if_already_exists", errParams)
	}
	if alreadyExists {
		return nil, gerrors.AlreadyExists("failed_to_create_trade.already_exists", errParams)
	}

	// Create trade.
	if err := dao.CreateTradeStrategy(ctx, trade); err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.dao", errParams)
	}

	// Read trade back out; we don't know the internal uuid, so we use our idempotency key which
	// guranteed to be unique.
	embelishedTrade, err := dao.ReadTradeStrategyByIdempotencyKey(ctx, trade.IdempotencyKey)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_create_trade.failed_to_read_created_trade_back", errParams)
	}

	return &tradeengineproto.CreateTradeStrategyResponse{
		TradeStrategyId: embelishedTrade.TradeStrategyID,
		Created:         timestamppb.New(embelishedTrade.Created),
	}, nil
}

func validateTradeStrategy(tradeStrategy *tradeengineproto.TradeStrategy) error {
	switch {
	case tradeStrategy == nil:
		return gerrors.BadParam("missing_param.tradeStrategy", nil)
	case tradeStrategy.Asset == "":
		return gerrors.BadParam("missing_param.asset", nil)
	case tradeStrategy.IdempotencyKey == "":
		return gerrors.BadParam("missing_param.idempotency_key", nil)
	case tradeStrategy.ActorId == "":
		return gerrors.BadParam("missing_param.actor_id", nil)
	case len(tradeStrategy.Entries) == 0:
		return gerrors.BadParam("missing_param.entries", nil)
	case tradeStrategy.InstrumentType == tradeengineproto.INSTRUMENT_TYPE_FUTURE_PERPETUAL && tradeStrategy.StopLoss == 0:
		return gerrors.FailedPrecondition("missing_param.stoploss_cannot_be_zero_for_futures_perpetuals", nil)
	}

	for _, entry := range tradeStrategy.Entries {
		if entry == 0 {
			return gerrors.BadParam("bad_param.zero_valued_entry", nil)
		}
	}

	return nil
}
