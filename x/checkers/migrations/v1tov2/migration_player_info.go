package v1tov2

import (
	"time"

	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func getOrNewPlayerInfoInMap(infoSoFar *map[string]*types.PlayerInfo, playerIndex string) (playerInfo *types.PlayerInfo) {
	playerInfo, found := (*infoSoFar)[playerIndex]
	if !found {
		playerInfo = &types.PlayerInfo{
			Index:          playerIndex,
			WonCount:       0,
			LostCount:      0,
			ForfeitedCount: 0,
		}
		(*infoSoFar)[playerIndex] = playerInfo
	}
	return playerInfo
}

func handleStoredGameChannel(ctx sdk.Context,
	k keeper.Keeper,
	gamesChannel <-chan []types.StoredGame,
	playerInfoChannel chan<- *types.PlayerInfo) {
	for games := range gamesChannel {
		playerInfos := make(map[string]*types.PlayerInfo, len(games))
		for _, game := range games {
			var winner string
			var loser string
			if game.Winner == rules.PieceStrings[rules.BLACK_PLAYER] {
				winner = game.Black
				loser = game.Red
			} else if game.Winner == rules.PieceStrings[rules.RED_PLAYER] {
				winner = game.Red
				loser = game.Black
			} else {
				continue
			}
			getOrNewPlayerInfoInMap(&playerInfos, winner).WonCount++
			getOrNewPlayerInfoInMap(&playerInfos, loser).LostCount++
		}
		for _, playerInfo := range playerInfos {
			playerInfoChannel <- playerInfo
		}
	}
	close(playerInfoChannel)
}

func handlePlayerInfoChannel(ctx sdk.Context, k keeper.Keeper,
	playerInfoChannel <-chan *types.PlayerInfo,
	done chan<- bool) {
	for receivedInfo := range playerInfoChannel {
		if receivedInfo != nil {
			existingInfo, found := k.GetPlayerInfo(ctx, receivedInfo.Index)
			if found {
				existingInfo.WonCount += receivedInfo.WonCount
				existingInfo.LostCount += receivedInfo.LostCount
				existingInfo.ForfeitedCount += receivedInfo.ForfeitedCount
			} else {
				existingInfo = *receivedInfo
			}
			k.SetPlayerInfo(ctx, existingInfo)
		}
	}
	done <- true
}

func MapStoredGamesReduceToPlayerInfo(ctx sdk.Context, k keeper.Keeper, chunk uint64) error {
	context := sdk.WrapSDKContext(ctx)
	response, err := k.StoredGameAll(context, &types.QueryAllStoredGameRequest{
		Pagination: &query.PageRequest{
			Limit: chunk,
		},
	})
	if err != nil {
		return err
	}
	gamesChannel := make(chan []types.StoredGame)
	playerInfoChannel := make(chan *types.PlayerInfo)
	done := make(chan bool)

	go handleStoredGameChannel(ctx, k, gamesChannel, playerInfoChannel)
	go handlePlayerInfoChannel(ctx, k, playerInfoChannel, done)
	gamesChannel <- response.StoredGame

	for response.Pagination.NextKey != nil {
		response, err = k.StoredGameAll(context, &types.QueryAllStoredGameRequest{
			Pagination: &query.PageRequest{
				Key:   response.Pagination.NextKey,
				Limit: chunk,
			},
		})
		if err != nil {
			return err
		}
		gamesChannel <- response.StoredGame
	}
	close(gamesChannel)

	<-done
	return nil
}

func addParsedCandidatesAndSort(parsedWinners []types.WinningPlayerParsed, candidates []types.WinningPlayerParsed) []types.WinningPlayerParsed {
	updated := append(parsedWinners, candidates...)
	types.SortWinners(updated)
	if types.LeaderboardWinnerLength < uint64(len(updated)) {
		updated = updated[:types.LeaderboardWinnerLength]
	}
	return updated
}

func AddCandidatesAndSortAtNow(parsedWinners []types.WinningPlayerParsed, now time.Time, playerInfos []types.PlayerInfo) []types.WinningPlayerParsed {
	parsedPlayers := make([]types.WinningPlayerParsed, 0, len(playerInfos))
	for _, playerInfo := range playerInfos {
		if playerInfo.WonCount > 0 {
			parsedPlayers = append(parsedPlayers, types.WinningPlayerParsed{
				PlayerAddress: playerInfo.Index,
				WonCount:      playerInfo.WonCount,
				DateAdded:     now,
			})
		}
	}
	return addParsedCandidatesAndSort(parsedWinners, parsedPlayers)
}

func AddCandidatesAndSort(parsedWinners []types.WinningPlayerParsed, ctx sdk.Context, playerInfos []types.PlayerInfo) []types.WinningPlayerParsed {
	return AddCandidatesAndSortAtNow(parsedWinners, types.GetDateAdded(ctx), playerInfos)
}

func handlePlayerInfosChannel(ctx sdk.Context, k keeper.Keeper,
	playerInfosChannel <-chan []types.PlayerInfo,
	done chan<- bool,
	chunk uint64) {
	winners := make([]types.WinningPlayerParsed, 0, types.LeaderboardWinnerLength+chunk)
	for receivedInfo := range playerInfosChannel {
		if receivedInfo != nil {
			winners = AddCandidatesAndSort(winners, ctx, receivedInfo)
		}
	}
	k.SetLeaderboard(ctx, types.CreateLeaderboardFromParsedWinners(winners))
	done <- true
}

func MapPlayerInfosReduceToLeaderboard(ctx sdk.Context, k keeper.Keeper, chunk uint64) error {
	context := sdk.WrapSDKContext(ctx)
	response, err := k.PlayerInfoAll(context, &types.QueryAllPlayerInfoRequest{
		Pagination: &query.PageRequest{
			Limit: PlayerInfoChunkSize,
		},
	})
	if err != nil {
		return err
	}
	playerInfosChannel := make(chan []types.PlayerInfo)
	done := make(chan bool)

	go handlePlayerInfosChannel(ctx, k, playerInfosChannel, done, chunk)
	playerInfosChannel <- response.PlayerInfo

	for response.Pagination.NextKey != nil {
		response, err = k.PlayerInfoAll(context, &types.QueryAllPlayerInfoRequest{
			Pagination: &query.PageRequest{
				Key:   response.Pagination.NextKey,
				Limit: PlayerInfoChunkSize,
			},
		})
		if err != nil {
			return err
		}
		playerInfosChannel <- response.PlayerInfo
	}
	close(playerInfosChannel)

	<-done
	return nil
}
