package Lottery

import "time"

type LotteriesCollector interface {
	 GetLotteryById (id string) (*Lottery, error)
	 GetLotteryByRaffler (id int) ([]*Lottery, error)
	 CreateLottery (lottery *Lottery) (*Lottery, error)
	 UpdateLottery (lottery *Lottery) (*Lottery, error)
	 SetWinner (lottery *Lottery) (*Lottery, error)
  	 AppendPlayer (lotID string,  playerID int, username string) error
  	 DeleteLottery  (ID string) error
	 GetLotteries () ([]*Lottery, error)

}

type Lottery struct {
	LotteryID string `json:"lottery_id"`
	LotName string `json:"lot_name"`
	Raffler User `json:"raffler"`
	Winner *User `json:"winner, omitempty"`
	Starttime time.Time `json:"startime"`
	Prize string `json:"prize"`
	TextMessage string `json:"text_message"`
	PlayerIDs []User `json:"player_ids"`
}

type User struct {
	Username string `json:"username"`
	UserID int `json:"user_id,omitempty"`
}

type MongoConfig struct {
	Host string
	Database string
	Port string
}
type ErrorResponse struct {
	StatusCode int
	ErrorMessage string
}

 func (er ErrorResponse) Error() string {
	return er.ErrorMessage
}