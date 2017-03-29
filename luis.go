package main

import (
	luis "github.com/kkdai/luis"
)

//LuisAction :
type LuisAction struct {
	LuisAPI *luis.Luis
}

//NewLuisAction :
func NewLuisAction(api string, appid string) *LuisAction {
	l := luis.NewLuis(api, appid)
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
