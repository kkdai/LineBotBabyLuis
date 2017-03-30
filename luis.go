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
	res, err := l.LuisAPI.IntentList()
	if err != nil {
		return nil, err
	}
	result := luis.NewIntentListResponse(res)
	return result, nil
}

//AddUtterance :
func (l *LuisAction) AddUtterance(intent, utterance string) {
	ex := luis.ExampleJson{ExampleText: utterance, SelectedIntentName: intent}
	res, err := l.LuisAPI.AddLabel(ex)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Done AddUtterance:", string(res))
}

//Predict :
func (l *LuisAction) Predict(utterance string) string {
	//Predict it, once you have train your models.
	res, err := l.LuisAPI.Predict(utterance)

	if err != nil {
		log.Println("Error happen on Predict:", err.Err)
		return ""
	}
	log.Println("Got response:", string(res))
	bestResult := luis.GetBestScoreIntent(luis.NewPredictResponse(res))
	fmt.Println("Get the best predict result:", bestResult)
	return bestResult.Name
}

//Train :
func (l *LuisAction) Train() {
	//Train it
	res, err := l.LuisAPI.Train()
	if err != nil {
		log.Println("Error happen on Trainning:", err.Err)
		return
	}
	log.Println("Training ret:", string(res))
}

// APPID = os.Getenv("APP_ID")
// 	API_KEY = os.Getenv("API_KEY")

// 	if API_KEY == "" {
// 		fmt.Println("Please export your key to environment first, `export SUB_KEY=12234 && export APP_ID=5678`")

// 	if API_KEY == "" {
// 		return
// 	}

// 	e := getLuis(t)

// 	res, err := e.IntelList()

// 	if err != nil {
// 		log.Error("Error happen on :", err.Err)
// 	}
// 	fmt.Println("Got response:", string(res))
// 	result := NewIntentListResponse(res)
// 	fmt.Println("Luis Intent Ret", result)

// 	//Add utterances
// 	ex := ExampleJson{ExampleText: "test", SelectedIntentName: "test2"}
// 	res, err = e.AddLabel(ex)

// 	//Train it
// 	res, err = e.Train()

// 	//Predict it, once you have train your models.
// 		res, err := e.Predict("test string")

// 	if err != nil {
// 		log.Error("Error happen on :", err.Err)
// 	}
// 	fmt.Println("Got response:", string(res))
// 	fmt.Println("Get the best predict result:", GetBestScoreIntent(NewPredictResponse(res)))
