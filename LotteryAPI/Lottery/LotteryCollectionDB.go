package Lottery

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LotteryCollection struct {
	collection *mongo.Collection
	dbcon *mongo.Database
}


func InitLotteryCollection (config MongoConfig) (LotteriesCollector, error) {
	clientOptions:=options.Client().ApplyURI("mongodb://"+config.Host+":"+config.Port)
	client,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		return nil,err
	}
	err = client.Ping(context.TODO(),nil)
	if err!=nil{
		return nil,err
	}

	db:=client.Database(config.Database)
	//collection=db.Collection("Lottery")
	return &LotteryCollection{collection:db.Collection("Lottery"),  dbcon:db,},nil
}


func (lc LotteryCollection) GetLotteryById(id string) (*Lottery, error) {
	filter:=bson.D{{"lotteryid",id}}
	lottery := &Lottery{}
	err := lc.collection.FindOne(context.TODO(),filter).Decode(&lottery)
	if err!=nil{
		return nil,err
	}
	return lottery, nil
}

func (lc LotteryCollection) GetLotteryByRaffler(id int) ([]*Lottery, error) {
	filter := bson.D{{"raffler.userid", id}}
	//findOptions:=options.Find()
	var lotteries []*Lottery
	cur,err := lc.collection.Find(context.TODO(), filter)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var lottery Lottery
		err:=cur.Decode(&lottery)
		if err!=nil{
			return nil,err
		}
		lotteries = append(lotteries ,&lottery)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return lotteries, nil




}

func (lc LotteryCollection) CreateLottery(lottery *Lottery) (*Lottery, error) {
	id := xid.New()
	lottery.LotteryID = id.String()
	lottery.PlayerIDs = make([]User, 0)
	insertResult,err := lc.collection.InsertOne(context.TODO(), lottery)
	if err!=nil{
		return nil,err
	}
	fmt.Println("Inserted document", insertResult.InsertedID)
	return lottery, nil
}

func (lc LotteryCollection) UpdateLottery(lottery *Lottery) (*Lottery, error) {

	filter:=bson.D{{"lotteryid", lottery.LotteryID}}

	update:=bson.D{{"$set",bson.D{
		{"lotname", lottery.LotName},
		{"prize",lottery.Prize},
		{"textmessage", lottery.TextMessage},
	}}}
	_,err := lc.collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return lottery,nil
}

func (lc LotteryCollection) SetWinner (lottery *Lottery) (*Lottery, error) {
	filter:=bson.D{{"lotteryid", lottery.LotteryID}}

	update:=bson.D{{"$set",bson.D{
		{"winner", bson.M{ "userid": lottery.Winner.UserID, "username": lottery.Winner.Username}},
	}}}
	_,err := lc.collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return lottery,nil
}



func (lc LotteryCollection) AppendPlayer(lotID string,  playerID int, username string) error {
	filter:=bson.D{{"lotteryid", lotID}}
	update:=bson.D{{"$push",bson.D{
		{"playerids", bson.M{ "userid": playerID, "username": username}},

	}}}
	_,err := lc.collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil
	}
	return nil
}

func (lc LotteryCollection) DeleteLottery(ID string) error {
	filter:=bson.D{{"lotteryid", ID}}
	_,err := lc.collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}

func (lc LotteryCollection) GetLotteries() ([]*Lottery, error) {
	findOptions:=options.Find()
	var lotteries []*Lottery
	cur,err := lc.collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var lottery Lottery
		err:=cur.Decode(&lottery)
		if err!=nil{
			return nil,err
		}
		lotteries = append(lotteries ,&lottery)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return lotteries, nil
}
