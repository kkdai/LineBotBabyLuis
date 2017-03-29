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

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var luisAction *LuisAction

func main() {
	var err error
	APPID := os.Getenv("APP_ID")
	API_KEY := os.Getenv("SUB_KEY")
	luisAction = NewLuisAction(APPID, API_KEY)

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
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
					log.Print(err)
				}

				ListAllIntents(bot, event.ReplyToken)
			}
		}
	}
}

func ListAllIntents(bot *linebot.Client, replyToken string) error {

	template := linebot.NewButtonsTemplate(
		"", "My button sample", "Hello, my button",
		linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
		linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
		linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
		linebot.NewMessageTemplateAction("Say message", "Rice=米"),
	)
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		return err
	}
	return nil
}
