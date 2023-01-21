package keeper

import (
	"fmt"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ForfeitExpiredGames(ctx sdk.Context) {
	opponents := map[string]string{
		rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER],
		rules.PieceStrings[rules.RED_PLAYER]:   rules.PieceStrings[rules.BLACK_PLAYER],
	}
	systemInfo, found := k.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	gameIndex := systemInfo.FifoHeadIndex
	var storedGame types.StoredGame

	for {
		if gameIndex == types.NoFifoIndex {
			break
		}
		storedGame, found = k.GetStoredGame(ctx, gameIndex)
		if !found {
			panic("Fifo head game not found " + systemInfo.FifoHeadIndex)
		}
		deadline, err := storedGame.GetDeadlineAsTime()
		if err != nil {
			panic(err)
		}

		if deadline.Before(ctx.BlockTime()) {
			// remove
			k.RemoveFromFifo(ctx, &storedGame, &systemInfo)

			lastBoard := storedGame.Board
			if storedGame.MoveCount <= 1 {
				// 没有真正开始的游戏，可以删掉
				k.RemoveStoredGame(ctx, gameIndex)

				if storedGame.MoveCount == 1 {
					k.MustRefundWager(ctx, &storedGame)
				}
			} else {
				storedGame.Winner, found = opponents[storedGame.Turn]
				if !found {
					panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
				}
				k.MustPayWinnings(ctx, &storedGame)

				// Here you can register a draw
				k.MustRegisterPlayerForfeit(ctx, &storedGame)

				storedGame.Board = ""
				k.SetStoredGame(ctx, storedGame)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.GameForfeitedEventType,
					sdk.NewAttribute(types.GameForfeitedEventGameIndex, gameIndex),
					sdk.NewAttribute(types.GameForfeitedEventWinner, storedGame.Winner),
					sdk.NewAttribute(types.GameForfeitedEventBoard, lastBoard),
				),
			)
			gameIndex = systemInfo.FifoHeadIndex
		} else {
			break
		}
	}

	k.SetSystemInfo(ctx, systemInfo)
}
