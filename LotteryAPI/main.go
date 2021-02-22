package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/mukhametkaly/OneLotteryAPI/LotteryAPI/Lottery"
)

func main()  {
	DBconfig := Lottery.MongoConfig{
		Host:     "172.19.0.1",
		Database: "test",
		Port:     "27017",
	}
	lotteryCollection, err := Lottery.InitLotteryCollection(DBconfig)
	if err != nil {
		panic(err)
	}
	lotteryExecute := Lottery.NewLotteryExecuter(lotteryCollection)
	lotteryEndpoint := Lottery.NewEndpointsFactory(lotteryExecute)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Methods("POST").Path("/lottery/create").HandlerFunc(lotteryEndpoint.CreateLottery())
	router.Methods("PUT").Path("/lottery/update").HandlerFunc(lotteryEndpoint.UpdateLottery())
	router.Methods("DELETE").Path("/lottery/delete/{id}").HandlerFunc(lotteryEndpoint.DeleteLottery("id"))
	router.Methods("GET").Path("/lottery/{lotid}/newplayer/{playerid}/username/{username}").HandlerFunc(lotteryEndpoint.AppendPlayer("lotid", "playerid", "username"))
	router.Methods("GET").Path("/lottery/{id}").HandlerFunc(lotteryEndpoint.GetLotteryById("id"))
	router.Methods("GET").Path("/lottery/raffler/{id}").HandlerFunc(lotteryEndpoint.GetLotteryByRaffler("id"))
	router.Methods("GET").Path("/lottery").HandlerFunc(lotteryEndpoint.GetLotteries())
	router.Methods("GET").Path("/lottery/play/{id}").HandlerFunc(lotteryEndpoint.PlayLottery("id"))
	fmt.Println("Server is running")
	http.ListenAndServe(":8000", router)
}