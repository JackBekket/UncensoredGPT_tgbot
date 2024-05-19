package command

import (
	"fmt"

	"github.com/JackBekket/uncensoredgpt_tgbot/lib/embeddings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/**
// Adds a new user to the database and assigns "Dialog_status" = 0.
func (c *Commander) AddNewUserToMap(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	c.usersDb[chatID] = database.User{
		ID:           chatID,
		Username:     updateMessage.From.UserName,
		DialogStatus: 0,
		Admin:        false,
	}

	user := c.usersDb[chatID]
	log.Printf(
		"Add new user to database: id: %v, username: %s\n",
		user.ID,
		user.Username,
	)

	msg := tgbotapi.NewMessage(user.ID, msgTemplates["hello"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Start!")),
	)
	c.bot.Send(msg)

	// check for registration
		registred := IsAlreadyRegistred(session, chatID)

		if registred {
			c.usersDb[chatID] = db.User{updateMessage.Chat.ID, updateMessage.Chat.UserName, 1}
		}



}
*/


func (c *Commander) HelpCommandMessage(updateMessage *tgbotapi.Message)  {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]
	msg := tgbotapi.NewMessage(user.ID, msgTemplates["help_command"])
	c.bot.Send(msg)
}

func (c *Commander) SearchDocuments(chatID int64, promt string, maxResults int) {
	//chatID := updateMessage.From.ID
	user := c.usersDb[chatID]
	results, err := embeddings.SemanticSearch(promt,maxResults)
	if err != nil {
		//return nil, err
		msg := tgbotapi.NewMessage(user.ID, "error occured: " + err.Error())
		c.bot.Send(msg)
	}

	for i, result := range results {
		content := result.PageContent
		msg := tgbotapi.NewMessage(user.ID, "result number: " + fmt.Sprint(i))
		c.bot.Send(msg)
		msg = tgbotapi.NewMessage(user.ID, "page content: " + content)
		c.bot.Send(msg)

		score := result.Score
		score_string := fmt.Sprintf("%f", score)

		msg = tgbotapi.NewMessage(user.ID, "score: " + score_string)
		c.bot.Send(msg)
	}

}

// Retrival-Augmented Generation
func (c *Commander) RAG(chatID int64, promt string, maxResults int) {
	user := c.usersDb[chatID]

	result, err := embeddings.RagSearch(promt,1)
	if err != nil {
		msg := tgbotapi.NewMessage(user.ID, "error occured: " + err.Error())
		c.bot.Send(msg)
	}
	msg := tgbotapi.NewMessage(user.ID, result)
	c.bot.Send(msg)
}

