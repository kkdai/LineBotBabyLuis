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

	luis "github.com/kkdai/luis"
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
	// luisAction = NewLuisAction(appID, apiKey)
	// res, err2 := luisAction.LuisAPI.IntentList()
	l := luis.NewLuis(apiKey, appID)
	res, err2 := l.IntentList()

	log.Println(res, err2)

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
			switch event.Message.(type) {
			case *linebot.TextMessage:
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				// 	log.Print(err)
				// }

				res, err := luisAction.GetIntents()
				if err != nil {
					log.Println(err)
				}

				var intentList []string
				log.Println("All intent:", *res)
				for _, v := range *res {
					intentList = append(intentList, v.Name)
				}

				ListAllIntents(bot, event.ReplyToken, intentList)
			}
		}
	}
}

//ListAllIntents :
func ListAllIntents(bot *linebot.Client, replyToken string, intents []string) {

	// var buttons []linebot.TemplateAction

	// for _, v := range intents {
	// 	buttons = append(buttons, linebot.NewPostbackTemplateAction(v, v, ""))
	// }

	// linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
	// linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),

	template := linebot.NewButtonsTemplate("", "My button sample", "Hello, my button",
		linebot.NewPostbackTemplateAction(intents[0], intents[0], ""),
		linebot.NewPostbackTemplateAction(intents[1], intents[1], ""),
		linebot.NewPostbackTemplateAction(intents[2], intents[2], ""),
		linebot.NewPostbackTemplateAction(intents[3], intents[3], ""))

	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		log.Print(err)
	}
}
