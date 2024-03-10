//package langchain

package main

import (
	"context"
	"log"
	"sync"

	db "github.com/JackBekket/uncensoredgpt_tgbot/lib/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mu = sync.Mutex{}

func SetupSequenceWithKey(
	bot *tgbotapi.BotAPI,
	user db.User,
	language string,
	ctx context.Context,
	spwd string,
	ai_endpoint string,
) {
	mu.Lock()
	defer mu.Unlock()
	chatID := user.ID
	gptKey := user.AiSession.GptKey
	log.Println("user GPT key from session: ", gptKey)
	//u_network := user.Network
	//log.Println("user network from session: ", u_network)
	log.Println("user model from session: ", user.AiSession.GptModel)

	u_pwd := user.AiSession.GptKey
	log.Println("upwd: ", u_pwd)

	// Initializing empty dialog thread
	thread, err := InitializeNewChatWithContextNoLimit(gptKey,user.AiSession.GptModel,ai_endpoint)
	if err != nil {
		log.Println(err)
	}
	user.AiSession.DialogThread = *thread
	db.UsersMap[chatID] = user // we need to store empty buffer *before* starting dialog

	

	switch language {
	case "English":
		probe, err := tryLanguage(user, "", 1, ctx,ai_endpoint,spwd,u_pwd)
		if err != nil {
			errorMessage(err, bot, user)
		} else {
			msg := tgbotapi.NewMessage(chatID, probe)
			bot.Send(msg)
			user.DialogStatus = 6
			db.UsersMap[chatID] = user
		}
	case "Russian":
		probe, err := tryLanguage(user, "", 2, ctx,ai_endpoint,spwd,u_pwd)
		if err != nil {
			errorMessage(err, bot, user)
		} else {
			msg := tgbotapi.NewMessage(chatID, probe)
			bot.Send(msg)
			user.DialogStatus = 6
			db.UsersMap[chatID] = user
		}
	default:
		probe, err := tryLanguage(user, language, 0, ctx,ai_endpoint,spwd,u_pwd)
		if err != nil {
			errorMessage(err, bot, user)
		} else {
			msg := tgbotapi.NewMessage(chatID, probe)
			bot.Send(msg)
			user.DialogStatus = 6
			db.UsersMap[chatID] = user
		}
	}

  
}

// LanguageCode: 0 - default, 1 - Russian, 2 - English
func tryLanguage(user db.User, language string, languageCode int, ctx context.Context, ai_endpoint string, spwd string, upwd string) (string, error) {
	var languagePromt string

	switch languageCode {
	case 1:
		languagePromt = "Hi, Do you speak english?"
	case 2:
		languagePromt = "Привет, ты говоришь по-русски?"
	default:
		languagePromt = language
	}

	log.Printf("Language: %v\n", languagePromt)
	//model := user.AiSession.GptModel
	thread := user.AiSession.DialogThread

	resp, err := ContinueChatWithContextNoLimit(&thread,languagePromt)
	if err != nil {
		return "", err
	} else {
		return resp, nil
	}

	/*
	resp, err := GenerateContentLAI(ai_endpoint,model,languagePromt)
	if err != nil {
		return "", err
	} else {
		LogResponseContentChoice(resp)
		answer := resp.Choices[0].Content
		return answer, nil
	}
	*/
}