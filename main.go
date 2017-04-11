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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	luis "github.com/kkdai/luis"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var luisAction *LuisAction
var allIntents *luis.IntentListResponse
var currentUtterance string

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
			case *linebot.ImageMessage:
				HandleHeavyContent(message.ID)
				// // content, err := bot.GetMessageContent(message.ID).Do()
				// // if err != nil {
				// // 	log.Println("imageMsg err:", err)
				// // }
				// // defer content.Content.Close()
				// log.Println("Img content:", content)

			case *linebot.TextMessage:
				ret := luisAction.Predict(message.Text)

				if ret.Name == "None" || ret.Name == "" || ret.Score < 0.5 {

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
					//List all intents
					ListAllIntents(bot, event.ReplyToken, intentList, message.Text)

				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("Daddy/Mommy, I just want to %s (%d %%)", ret.Name, int(ret.Score*100)))).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		} else if event.Type == linebot.EventTypePostback {
			//Add new utterance into original intent
			luisAction.AddUtterance(event.Postback.Data, currentUtterance)

			retStr := fmt.Sprintf("Daddy/Mommy, I just learn new utterance %s for intent: %s.", currentUtterance, event.Postback.Data)
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(retStr)).Do(); err != nil {
				log.Print(err)
			}

			//Train it right away
			luisAction.Train()
		}
	}
}

//ListAllIntents :
func ListAllIntents(bot *linebot.Client, replyToken string, intents []string, utterance string) {
	askStmt := fmt.Sprintf("Your utterance %s is not exist, please select correct intent.", utterance)
	log.Println("askStmt:", askStmt)

	var sliceTemplateAction []linebot.TemplateAction
	for _, v := range intents {
		sliceTemplateAction = append(sliceTemplateAction, linebot.NewPostbackTemplateAction(v, v, ""))
	}

	template := linebot.NewButtonsTemplate("", "Select your intent for your baby", utterance, sliceTemplateAction...)

	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTemplateMessage("Select your intent for your baby", template)).Do(); err != nil {
		log.Print(err)
	}
	currentUtterance = utterance
}

//HandleHeavyContent :
func HandleHeavyContent(messageID string) error {
	content, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		log.Println("Get msg err:", err)
		return err
	}
	defer content.Content.Close()
	log.Printf("Got file: %s", content.ContentType)
	_, err = SaveContent(content.Content)
	if err != nil {
		log.Println("Save file err:", err)
		return err
	}
	return nil
}

//HandleImage :
func HandleImage(message *linebot.ImageMessage, replyToken string) error {
	return HandleHeavyContent(message.ID)
}

// SaveContent :
func SaveContent(content io.ReadCloser) (*os.File, error) {

	file, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return nil, err
	}
	log.Printf("Saved to %s", file.Name())

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("image", file.Name())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Println(err)
		return nil, err
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	if err = w.Close(); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Total file length:", len(b))
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", "http://107.167.183.27:3000/api/v1/tf-image/", &b)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	log.Println("data content type:", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		log.Println(err)
	}

	// res, err := http.Post("http://107.167.183.27:3000/api/v1/tf-image/", "binary/octet-stream", file)
	// if err != nil {
	// 	panic(err)
	// }
	// defer res.Body.Close()
	// message, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("Ret:%s", string(message))

	return file, nil
}
