package command

import (
	"log"
	"net/url"
	"path"

	"github.com/JackBekket/uncensoredgpt_tgbot/lib/langchain"
	"github.com/JackBekket/uncensoredgpt_tgbot/lib/localai"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)



type WrongPwdError struct {
	message string
}

// Message:	case0 - "Input your openAI API key. It can be created at https://platform.openai.com/accousernamet/api-keys".
//  DialogStatus 2 -> 3
func (c *Commander) InputYourAPIKey(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(
		user.ID,
		msgTemplates["case0"],
	)
	c.bot.Send(msg)

	user.DialogStatus = 3
	c.usersDb[chatID] = user
}


// DialogStatus 0 - > 1
func (c *Commander) ChooseNetwork(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]
	// render menu
	msg := tgbotapi.NewMessage(user.ID, msgTemplates["ch_network"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("openai"),
			tgbotapi.NewKeyboardButton("localai")),

	)
	c.bot.Send(msg)

	user.DialogStatus = 1	  // this is output dialog status
	c.usersDb[chatID] = user // commit changes

}


// Dialog status 1 -> 2
func (c *Commander) HandleNetworkChoose(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	network := updateMessage.Text
	user := c.usersDb[chatID]
	switch network {
	case "openai":

		user.Network = network
		user.AiSession.AI_Type = 0
		user.DialogStatus = 2
		c.usersDb[chatID] = user
		c.InputYourAPIKey(updateMessage)
	case "localai":

		user.Network = network
		user.AiSession.AI_Type = 1
		user.DialogStatus = 2
		c.usersDb[chatID] = user
		c.InputYourAPIKey(updateMessage)
	default:
		c.WrongNetwork(updateMessage)
	}

}




//	update Dialog_Status 3 -> 4
func (c *Commander) ChooseModel(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	gptKey := updateMessage.Text	// handling previouse message
	user := c.usersDb[chatID]
	network := user.Network

	// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
	// Since this part is oftenly get an usernamecaught exeption, we debug what user input as key. It's bad, I know, but usernametil we got key validation we need this part.
	log.Println("Key promt: ", gptKey)
	user.AiSession.GptKey = gptKey // store key in memory

	switch network {
	case "localai" :
		c.RenderModelMenuLAI(chatID)
		user.DialogStatus = 4
		c.usersDb[chatID] = user
	
	case "openai" :
		c.RenderModelMenuOAI(chatID)
		user.DialogStatus = 4
		c.usersDb[chatID] = user


	default :
		c.WrongNetwork(updateMessage)
	}
}


// DialogStatus 4 -> 5
func (c *Commander) HandleModelChoose(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	model_name := updateMessage.Text
	user := c.usersDb[chatID]
	network := user.Network

	switch network {
	case "localai" :
		switch model_name {
		case "wizard-uncensored-13b":
			c.attachModel(model_name, chatID)
			user.AiSession.GptModel = model_name
			c.RenderLanguage(chatID)
	
			user.DialogStatus = 5
			c.usersDb[chatID] = user
		case "wizard-uncensored-30b":
			c.attachModel(model_name, chatID)
			user.AiSession.GptModel = model_name
			c.RenderLanguage(chatID)
	
			user.DialogStatus = 5
			c.usersDb[chatID] = user
		default:
			c.WrongModel(updateMessage)
		}

	case "openai" :
		switch model_name {
		case "gpt-3.5":
			model_name = "gpt-3.5-turbo"
			c.attachModel(model_name, chatID)
			user.AiSession.GptModel = model_name
			c.RenderLanguage(chatID)
	
			user.DialogStatus = 5
			c.usersDb[chatID] = user
		case "gpt-4":
			c.attachModel(model_name, chatID)
			user.AiSession.GptModel = model_name
			c.RenderLanguage(chatID)
	
			user.DialogStatus = 5
			c.usersDb[chatID] = user
		default:
			c.WrongModel(updateMessage)
		}
	}


}



// Depracated?
// Message: "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language".
//
//	update Dialog_Status = 3
func (c *Commander) ModelGPT3DOT5(updateMessage *tgbotapi.Message) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", updateMessage.Text)

	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	//modelName := openai.GPT3Dot5Turbo // gpt-3.5
	modelName := "gpt-3.5"

	user.AiSession.GptModel = modelName
	msg := tgbotapi.NewMessage(user.ID, "your session model: "+modelName)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(user.ID, "Choose a language or send 'Hello' in your desired language.")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("English"),
			tgbotapi.NewKeyboardButton("Russian")),
	)
	c.bot.Send(msg)

	user.DialogStatus = 3
	c.usersDb[chatID] = user
}


// render language menu
func (c *Commander) RenderLanguage(chat_id int64) {
	chatID := chat_id
	//user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(chatID, "Choose a language or send 'Hello' in your desired language.")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("English"),
			tgbotapi.NewKeyboardButton("Russian")),
	)
	c.bot.Send(msg)

	//user.DialogStatus = 3
	//c.usersDb[chatID] = user
}

// low level attach model name to user profile
func (c *Commander) attachModel(model_name string, chatID int64) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", model_name)

	user := c.usersDb[chatID]

	modelName := model_name
	user.AiSession.GptModel = modelName
	msg := tgbotapi.NewMessage(user.ID, "your session model: "+modelName)
	c.bot.Send(msg)
	c.usersDb[chatID] = user
}

// internal for attach api key to a user
func (c *Commander) AttachKey(gpt_key string, chatID int64) {
	log.Println("Key promt: ", gpt_key)
	user := c.usersDb[chatID]
	user.AiSession.GptKey = gpt_key // store key in memory
	c.usersDb[chatID] = user
}

// Dangerouse! NOTE -- probably work only internal
func (c *Commander) ChangeDialogStatus(chatID int64, ds int8) {
	user := c.usersDb[chatID]
	old_status := user.DialogStatus
	log.Println("dialog status changed, old status is ", old_status)
	log.Println("new status is ", ds)
	user.DialogStatus = ds
}

func (c *Commander) RenderModelMenuOAI(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, msgTemplates["case1"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("gpt-3.5"),
		tgbotapi.NewKeyboardButton("gpt-4")),
	)
	c.bot.Send(msg)
}

func (c *Commander) RenderModelMenuLAI(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, msgTemplates["case1"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("wizard-uncensored-13b")),
	//	tgbotapi.NewKeyboardButton("wizard-uncensored-30b")),
	)
	c.bot.Send(msg)
}

// update Dialog_Status = 4
func (c *Commander) WrongModel(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(user.ID, "type wizard-uncensored-13b")
	c.bot.Send(msg)

	user.DialogStatus = 4
	c.usersDb[chatID] = user
}

// update Dialog_Status = 0
func (c *Commander) WrongNetwork(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(user.ID, "type openai or localai")
	c.bot.Send(msg)

	user.DialogStatus = 0
	c.usersDb[chatID] = user
}


// update update Dialog_Status 5 -> 6
func (c *Commander) ConnectingToAiWithLanguage(updateMessage *tgbotapi.Message, ai_endpoint string) {
	chatID := updateMessage.From.ID
	language := updateMessage.Text
	user := c.usersDb[chatID]
	log.Println("check gpt key exist:", user.AiSession.GptKey)

	network := user.Network

	msg := tgbotapi.NewMessage(user.ID, "connecting to ai node")
	c.bot.Send(msg)

	//go localai.SetupSequenceWithKey(c.bot, user, language, c.ctx, lpwd, ai_endpoint)

		if network == "localai" {
			go langchain.SetupSequenceWithKey(c.bot,user,language,c.ctx,ai_endpoint)
		} else {

		go langchain.SetupSequenceWithKey(c.bot,user,language,c.ctx,"")
		//go localai.SetupSequenceWithKey(c.bot,user,language,c.ctx,lpwd,ai_endpoint)
		}
	
	
}

// Generates an image with the /image command.
//
// Generates and sends text to the user. This is *main loop*
//
// update Dialog_Status 6 -> 6 (loop), 
func (c *Commander) DialogSequence(updateMessage *tgbotapi.Message, ai_endpoint string) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]
	switch updateMessage.Command() {
	
		case "image":
			msg := tgbotapi.NewMessage(user.ID, "Image link generation...")
			c.bot.Send(msg)

			promt := updateMessage.CommandArguments()
			log.Printf("Command /image arg: %s\n", promt)
			if (promt == "") {
				c.GenerateNewImageLAI_SD("evangelion, neon, anime",chatID,c.bot)
			} else {
				c.GenerateNewImageLAI_SD(promt,chatID,c.bot)
			}
			//go openaibot.StartImageSequence(c.bot, updateMessage, chatID, promt, c.ctx)

	
	default:
		promt := updateMessage.Text
		//go localai.StartDialogSequence(c.bot, chatID, promt, c.ctx, ai_endpoint)
		go langchain.StartDialogSequence(c.bot,chatID,promt,c.ctx,ai_endpoint)
	}	
}

// stable diffusion
func (c *Commander) GenerateNewImageLAI_SD(promt string, chatID int64, bot *tgbotapi.BotAPI) {
	size := "256x256"
	filepath, err := localai.GenerateImageStableDissusion(promt, size)
	if err != nil {
		//return nil, err
		log.Println(err)
	}
	log.Println("url_path: ", filepath)
	sendImage(bot, chatID, filepath)
}

func sendImage(bot *tgbotapi.BotAPI, chatID int64, path string) {
	// Prepare a photo message
	fileName := transformURL(path)
	log.Println("local file name: ", fileName)

	telegraphLink := localai.UploadToTelegraph(fileName)
	log.Println("uploaded to telegraph successfully, link is: ", telegraphLink)

	// Path to the image/file locally
	// filePath := "/path/to/image.png" + local_path
	/*
			 // Creating a LocalFile object from the local path
			photoBytes, err := ioutil.ReadFile(filePath)
			if err != nil {
		    	log.Println(err)
						}
			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: photoBytes,
				}
	*/
	//message, err := bot.Send(tgbotapi.NewPhotoUpload(int64(chatID), photoFileBytes))
	/* photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(local_path))
	if _, err := bot.Send(photo); err != nil {
	log.Fatalln(err)
	} */
	msg := tgbotapi.NewMessage(chatID, telegraphLink)
	bot.Send(msg)
}

func transformURL(inputURL string) string {
	// Replace "http://localhost:8080" with "/tmp" using strings.Replace
	parsedURL, _ := url.Parse(inputURL)

	// Use path.Base to get the filename from the URL path
	fileName := path.Base(parsedURL.Path)
	return fileName
}


func (c *Commander) CheckLocalPWD(upwd string, spwd string) (bool, error) {
	if upwd != spwd {
		err := &WrongPwdError{"wrong password"}
		return false, err
	} else {
		return true, nil
	}
}

func (e *WrongPwdError) Error() string {
    return e.message
}