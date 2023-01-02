package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sysInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("System info not found!")
	}
	newIndex := strconv.FormatUint(sysInfo.NextId, 10)

	newGame := rules.New()
	storedGame := types.StoredGame{
		Index:       newIndex,
		Board:       newGame.String(),
		Turn:        rules.PieceStrings[newGame.Turn],
		Black:       msg.Black,
		Red:         msg.Red,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner:      rules.PieceStrings[rules.NO_PLAYER],
		Wager:       msg.Wager,
	}

	if err := storedGame.Validate(); err != nil {
		return nil, err
	}

	k.Keeper.SendToFifoTail(ctx, &storedGame, &sysInfo)
	k.Keeper.SetStoredGame(ctx, storedGame)
	sysInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, sysInfo)
	ctx.GasMeter().ConsumeGas(types.CreateGameGas, "Create game")

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
			sdk.NewAttribute(types.GameCreatedEventWager, strconv.FormatUint(msg.Wager, 10)),
		),
	)
	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
