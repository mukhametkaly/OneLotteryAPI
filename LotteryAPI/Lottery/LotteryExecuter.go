package Lottery

import "time"

type LotteryExecuter interface {
	GetLotteryById (id string) (*Lottery, error)
	GetLotteryByRaffler (id int) ([]*Lottery, error)
	CreateLottery (lottery *Lottery) (*Lottery, error)
	UpdateLottery (lottery *Lottery) (*Lottery, error)
	AppendPlayer (lotID string,  playerID int, username string) error
	DeleteLottery  (ID string) error
	GetLotteries () ([]*Lottery, error)
	IsWinnerExist (id string) error
	Play(id string) (*Lottery, error)
}

type LotteryExecute struct {
	LotteryColl LotteriesCollector
}

func NewLotteryExecuter(lotColletion LotteriesCollector) LotteryExecuter {
	return &LotteryExecute{
		LotteryColl: lotColletion,
	}
}

func (l LotteryExecute) GetLotteryById(id string) (*Lottery, error) {

	data, err := l.LotteryColl.GetLotteryById(id)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return nil, errResp
	}
	if data == nil {
		errResp := ErrorResponse{
			StatusCode:   404,
			ErrorMessage: "Not Found",
		}
		return nil, errResp
	}

	return  data, nil

}

func (l LotteryExecute) GetLotteryByRaffler(id int) ([]*Lottery, error) {
	data, err := l.LotteryColl.GetLotteryByRaffler(id)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return nil, errResp
	}
	if data == nil {
		errResp := ErrorResponse{
			StatusCode:   404,
			ErrorMessage: "Not Found",
		}
		return nil, errResp
	}
	return data, nil


}

func (l LotteryExecute) CreateLottery(lottery *Lottery) (*Lottery, error) {
	if !lottery.Chek() {
		errResp := ErrorResponse{
			StatusCode:   400,
			ErrorMessage: "Bad request",
		}
		return nil, errResp
	}
	start := time.Now()

	lottery.Starttime = start

	response, err := l.LotteryColl.CreateLottery(lottery)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return nil, errResp
	}

	return response, err

}

func (l LotteryExecute) UpdateLottery(lottery *Lottery) (*Lottery, error) {

	if !lottery.Chek() {
		errResp := ErrorResponse{
			StatusCode:   400,
			ErrorMessage: "Bad request",
		}
		return nil, errResp
	}

	if err := l.IsWinnerExist(lottery.LotteryID); err != nil {
		return nil, err
	}

	updatedlottery, err := l.LotteryColl.UpdateLottery(lottery)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return nil, errResp
	}

	return updatedlottery, err

}

func (l LotteryExecute) AppendPlayer(lotID string, playerID int, username string) error {

	lottery, err := l.GetLotteryById(lotID)

	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return errResp
	}

	if lottery.Winner != nil {
		errResp := ErrorResponse{
			StatusCode:   405,
			ErrorMessage: "WinnerExist",
		}
		return errResp
	}

	for _, i := range lottery.PlayerIDs {
		if i.UserID == playerID {
			errResp := ErrorResponse{
				StatusCode:   405,
				ErrorMessage: "PlayerAlreadyEnjoy",
			}
			return errResp
		}
	}


	err = l.LotteryColl.AppendPlayer(lotID, playerID, username)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return errResp
	}

	return nil

}

func (l LotteryExecute) DeleteLottery(ID string) error {
	return l.LotteryColl.DeleteLottery(ID)
}

func (l LotteryExecute) GetLotteries() ([]*Lottery, error) {
	return l.LotteryColl.GetLotteries()
}

func (l LotteryExecute) IsWinnerExist(id string) error {

	lottery, err := l.GetLotteryById(id)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return errResp
	}

	if lottery.Winner != nil {
		errResp := ErrorResponse{
			StatusCode:   405,
			ErrorMessage: "WinnerExist",
		}
		return errResp
	}
	return nil
}

func (l LotteryExecute) Play(id string) (*Lottery, error) {

	lottery, err := l.GetLotteryById(id)
	if err != nil {
		return nil, err
	}

	if lottery.Winner != nil {
		errResp := ErrorResponse{
			StatusCode:   405,
			ErrorMessage: "WinnerExist",
		}
		return nil, errResp
	}

	if !lottery.Play() {
		errResp := ErrorResponse{
			StatusCode:   405,
			ErrorMessage: "NoPlayers",
		}
		return nil, errResp
	}

	LotWithWinner, err := l.LotteryColl.SetWinner(lottery)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		return nil, errResp
	}

	return LotWithWinner, nil



}

