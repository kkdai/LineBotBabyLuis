// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kkdai/luis"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var luisAction *LuisAction
var allIntents *luis.IntentListResponse

func main() {
	var err error
	appID := os.Getenv("APP_ID")
	apiKey := os.Getenv("APP_KEY")
	log.Println("Luis:", appID, apiKey)
	luisAction = NewLuisAction(appID, apiKey)

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				ret := luisAction.Predict(message.Text)
				if ret == "None" || ret == "" {

					res, err := luisAction.GetIntents()
					if err != nil {
						log.Println(err)
						return
					}
					var intentList []string
					log.Println("All intent:", *res)
					for _, v := range *res {
						if v.Name != "None" {
							intentList = append(intentList, v.Name)
						}
					}

					ListAllIntents(bot, event.ReplyToken, intentList, message.Text)

				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("Hi Dady/Mam: I just want to :%s", ret))).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

//ListAllIntents :
func ListAllIntents(bot *linebot.Client, replyToken string, intents []string, utterance string) {

	// var buttons []linebot.TemplateAction

	askStmt := fmt.Sprintf("Your utterance %s is not exist, please select correct intent.", utterance)
	log.Println("askStmt:", askStmt)

	// template := linebot.NewButtonsTemplate("", "Please select your intent of your word", "test",
	template := linebot.NewButtonsTemplate("", "Select your intent", askStmt,
		linebot.NewPostbackTemplateAction(intents[0], intents[0], ""),
		linebot.NewPostbackTemplateAction(intents[1], intents[1], ""),
		linebot.NewPostbackTemplateAction(intents[2], intents[2], ""),
		linebot.NewPostbackTemplateAction(intents[3], intents[3], ""))

	//	if _, err := bot.ReplyMessage(replyToken, linebot.NewTemplateMessage("test....", template)).Do(); err != nil {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTemplateMessage("Select your intent", template)).Do(); err != nil {
		log.Print(err)
	}
}
