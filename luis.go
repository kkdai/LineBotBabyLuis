package main

import (
	"fmt"
	"log"

	luis "github.com/kkdai/luis"
)

//LuisAction :
type LuisAction struct {
	LuisAPI *luis.Luis
}

//NewLuisAction :
func NewLuisAction(appid string, appkey string) *LuisAction {
	l := luis.NewLuis(appkey, appid)
	return &LuisAction{LuisAPI: l}
}

//GetIntents :
func (l *LuisAction) GetIntents() (*luis.IntentListResponse, *luis.ErrorResponse) {
	res, err := l.LuisAPI.IntelList()
	if err != nil {
		return nil, err
	}
	result := luis.NewIntentListResponse(res)
	return result, nil
}

//AddUtterance :Add new example to your intent.
func (l *LuisAction) AddUtterance(intent, utterance string) {
	ex := luis.ExampleJson{ExampleText: utterance, SelectedIntentName: intent}
	res, err := l.LuisAPI.AddLabel(ex)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Done AddUtterance:", string(res))
}

//Predict :Predict your example(utterance).
//Note: If your model never train before, use predict will get error.
func (l *LuisAction) Predict(utterance string) luis.IntentScore {
	//Predict it, once you have train your models.
	res, err := l.LuisAPI.Predict(utterance)

	if err != nil {
		log.Println("Error happen on Predict:", err.Err)
		return luis.IntentScore{}
	}
	log.Println("Got response:", string(res))
	bestResult := luis.GetBestScoreIntent(luis.NewPredictResponse(res))
	fmt.Println("Get the best predict result:", bestResult)
	return bestResult
}

//Train :Train your model, this is async call. It might take time if your data set is big.
func (l *LuisAction) Train() {
	//Train it
	res, err := l.LuisAPI.Train()
	if err != nil {
		log.Println("Error happen on Trainning:", err.Err)
		return
	}
	log.Println("Training ret:", string(res))
}
