package keeper

import (
	"context"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
	}

	switch msg.Creator {
	case storedGame.Black:
		if 0 < storedGame.MoveCount { // Notice the use of the new field
			return nil, types.ErrBlackAlreadyPlayed
		}
	case storedGame.Red:
		if 1 < storedGame.MoveCount { // Notice the use of the new field
			return nil, types.ErrRedAlreadyPlayed
		}
	default:
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	}

	k.Keeper.MustRefundWager(ctx, &storedGame)

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}
	k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)
	k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	refund := uint64(types.RejectGameRefundGas)
	if consumed := ctx.GasMeter().GasConsumed(); consumed < refund {
		refund = consumed
	}
	ctx.GasMeter().RefundGas(refund, "Reject game")

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.GameRejectedEventType,
		sdk.NewAttribute(types.GameRejectedEventCreator, msg.Creator),
		sdk.NewAttribute(types.GameRejectedEventGameIndex, msg.GameIndex),
	))
	return &types.MsgRejectGameResponse{}, nil
}
