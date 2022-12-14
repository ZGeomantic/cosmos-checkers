package keeper

import (
	"context"

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

	k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.GameRejectedEventType,
		sdk.NewAttribute(types.GameRejectedEventCreator, msg.Creator),
		sdk.NewAttribute(types.GameRejectedEventGameIndex, msg.GameIndex),
	))
	return &types.MsgRejectGameResponse{}, nil
}
