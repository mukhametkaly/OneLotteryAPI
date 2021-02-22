package Lottery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type LotteryEndpoints interface {
	GetLotteryById (idParam string) 									func(w http.ResponseWriter,r *http.Request)
	GetLotteryByRaffler (idParam string) 								func(w http.ResponseWriter,r *http.Request)
	CreateLottery ()													func(w http.ResponseWriter,r *http.Request)
	UpdateLottery () 													func(w http.ResponseWriter,r *http.Request)
	AppendPlayer  (lotID string,  playerID string, username string) 	func(w http.ResponseWriter,r *http.Request)
	DeleteLottery (idParam string) 										func(w http.ResponseWriter,r *http.Request)
	GetLotteries  () 													func(w http.ResponseWriter,r *http.Request)
	PlayLottery(idParam string) 										func(w http.ResponseWriter,r *http.Request)
}

type LotteryEndpointsFactory struct {
	LotteryExec LotteryExecuter
}


func NewEndpointsFactory(lotExec LotteryExecuter) LotteryEndpoints {
	return &LotteryEndpointsFactory{
		LotteryExec: lotExec,
	}
}

func encodeError(err error, w http.ResponseWriter) {
	servErr := err.(ErrorResponse)

	switch servErr.StatusCode {
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry :( \n Server Internal error"))
		fmt.Println(err.Error())
	case 400:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		fmt.Println(err.Error())
	case 404:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		fmt.Println(err.Error())
	case 413:
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("Request entity too large"))
		fmt.Println(err.Error())
	case 503:
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service unavailable"))
		fmt.Println(err.Error())
	case 405:
		if servErr.Error() == "WinnerExist" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Sorry, the error could be due to the fact that there is already a winner in this lottery"))
			fmt.Println(err.Error())
		} else if servErr.Error() == "NoPlayers" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Sorry not enough players in lottery"))
			fmt.Println(err.Error())
		} else if servErr.Error() == "PlayerAlreadyEnjoy" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("You are already enjoy"))
			fmt.Println(err.Error())
		}

	}


}

func respondJSON(w http.ResponseWriter, lottery *Lottery) {
	response, err := json.Marshal(*lottery)
	if err != nil {
		errResp := ErrorResponse{
			StatusCode:   500,
			ErrorMessage: err.Error(),
		}
		encodeError(errResp, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func (lef LotteryEndpointsFactory) GetLotteryById(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		paramid, paramerr:=vars[idParam]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}

		data, err := lef.LotteryExec.GetLotteryById(paramid)
		if err != nil {
			encodeError(err, w)
			return
		}
		respondJSON(w, data)
	}
}

func (lef LotteryEndpointsFactory) GetLotteryByRaffler(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		paramid, paramerr:=vars[idParam]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}
		id, err := strconv.Atoi(paramid)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}

		data, err := lef.LotteryExec.GetLotteryByRaffler(id)
		if err != nil {
			encodeError(err, w )
			return
		}

		response, err := json.Marshal(data)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		w.WriteHeader(200)
		w.Write(response)
	}
}

func (lef LotteryEndpointsFactory) CreateLottery() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		lottery := &Lottery{}
		err = json.Unmarshal(data, &lottery)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}

		lottery, err = lef.LotteryExec.CreateLottery(lottery)
		if err != nil {
			encodeError(err, w )
			return
		}
		respondJSON(w, lottery)
	}
}

func (lef LotteryEndpointsFactory) UpdateLottery() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}

		updatedLottery := &Lottery{}
		err = json.Unmarshal(data, &updatedLottery)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}

		updatedLottery, err = lef.LotteryExec.UpdateLottery(updatedLottery)
		if err != nil {
			encodeError(err, w )
			return
		}
		respondJSON(w, updatedLottery)
	}
}

func (lef LotteryEndpointsFactory) AppendPlayer(idParam, UserID, Username string ) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		lotID, paramerr:=vars[idParam]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}
		paramid, paramerr:=vars[UserID]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}
		userID, err := strconv.Atoi(paramid)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		username, paramerr:=vars[Username]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}

		err = lef.LotteryExec.AppendPlayer(lotID, userID, username)
		if err != nil {
			encodeError(err, w )
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("User enjoy"))
	}
}

func (lef LotteryEndpointsFactory) DeleteLottery(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:=mux.Vars(r)
		id, paramerr:=vars[idParam]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}
		err := lef.LotteryExec.DeleteLottery(id)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("Lot deleted"))
	}
}

func (lef LotteryEndpointsFactory) GetLotteries() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lotteries, err := lef.LotteryExec.GetLotteries()
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		response, err := json.Marshal(lotteries)
		if err != nil {
			errResp := ErrorResponse{
				StatusCode:   500,
				ErrorMessage: err.Error(),
			}
			encodeError(errResp, w )
			return
		}
		w.WriteHeader(200)
		w.Write(response)
	}
}



func (lef LotteryEndpointsFactory) PlayLottery(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars:=mux.Vars(r)
		id, paramerr:=vars[idParam]
		if !paramerr{
			errResp := ErrorResponse{
				StatusCode:   400,
				ErrorMessage: "No arguments",
			}
			encodeError(errResp, w )
			return
		}

		LotteryWithWinner, err := lef.LotteryExec.Play(id)
		if err != nil {
			encodeError(err, w )
			return
		}
		respondJSON(w, LotteryWithWinner)

	}

}