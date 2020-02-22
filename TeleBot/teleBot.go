package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("<token>")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		getFileFromTelegram(update)

	}
}
func getFileFromTelegram(update tgbotapi.Update) {
	fmt.Println(update.Message.Voice)
	if update.Message == nil { // ignore any non-Message Updates
		return
	}
	fmt.Println(update.Message.Voice)
	if update.Message.Voice == nil { // ignore if no Voice-Message Updates
		return
	}
	fileId := update.Message.Voice.FileID
	fmt.Println(fileId)
	if fileId != "" {

		fileUrl := "https://api.telegram.org/bot" + "<token>" + "/getFile?file_id=" + fileId
		resp, err := http.Get(fileUrl)
		if err != nil {
			return
		}
		filePath, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fPath := strings.Split(string(filePath), "\"")
		//fmt.Println(fPath[len(fPath)-2])
		fileUrl = "https://api.telegram.org/file/bot<token>/" + fPath[len(fPath)-2]
		fPath = strings.Split(fPath[len(fPath)-2], "/")
		if err := DownloadFile(fPath[len(fPath)-1], fileUrl); err != nil {
			panic(err)
		}
	}
}
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
