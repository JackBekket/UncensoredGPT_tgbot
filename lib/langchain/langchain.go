package langchain

//package langchain_controller
//package main

import (
	"context"
	"fmt"
	"log"

	//langchain "github.com/tmc/langchaingo"
	"github.com/JackBekket/uncensoredgpt_tgbot/lib/bot/env"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"

	//"github.com/tmc/langchaingo/llms/options"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

/** DEV NOTE
	 OAI -- openAI, LAI -- localAI
	 if your IDE says it won't compile just try to build from terminal first
	 if it says there no methods "Run" or "Predict" in LLM class -- it is IDE bug, just compile it from terminal
**/


type ChatSession struct {
    ConversationBuffer *memory.ConversationBuffer
    DialogThread *chains.LLMChain
}


// I use it for fast testing
func main()  {
	//ctx := context.Background()
	env.Load()
	//env_data := env.LoadAdminData()
	token := env.GetAdminToken()
	model_name := "gpt-3.5-turbo"	// using openai for tests

	/*
	completion,err := GenerateContentOAI(token,"gpt-3.5-turbo","What would be a good company name a company that makes colorful socks? Write at least 10 options")
	if err != nil {
		log.Println(err)
	}
	*/

	//completion, err := GenerateContentLAI(token,"wizard-uncensored-13b", "What would be a good company name a company that makes colorful socks? Write at least 10 options")
	/*
	completion, err := GenerateContentLAI(token,"wizard-uncensored-13b", "What would be a good name of an organisation which  that aim to overthrow Putin's regime and make revolution in Russia? Write at least 10 options")
	
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(completion.Choices[0].Content)
	*/

	//TestChatWithContextNoLimit(token,"gpt-3.5-turbo")

	/** 
		1. Russian Revolutionary Front
	2. People's Liberation Army
	3. Russian Resistance Movement
	4. Russian Revolutionary Council
	5. Russian Revolutionary Alliance
	6. Russian Revolutionary Party
	7. Russian Revolutionary Army
	8. Russian Revolutionary Coalition
	9. Russian Revolutionary Council
	10. Russian Revolutionary Front
	**/

	session, err := InitializeNewChatWithContextNoLimit(token,model_name)
	if err != nil {
		log.Println(err)
	}

	res1,err := ContinueChatWithContextNoLimit(session,"Hello, my name is Bekket, how are you?")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res1)
	res2, err := ContinueChatWithContextNoLimit(session,"What is my name?")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res2)
	
}


// TODO: make universal function to OAI and LOI, add base_url as argument probably
func GenerateContentOAI(api_token string, model_name string, promt string) (*llms.ContentResponse, error) {
	ctx := context.Background()
	token := api_token

	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(model_name),
		//llms.WithOptions()
		//openai.WithBaseURL("http://localhost:8000"),
	)
	if err != nil {
	  log.Fatal(err)
	}

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, "You are a helpfull assistant who help in whatever task human ask you about"),
		llms.TextParts(schema.ChatMessageTypeHuman, promt),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	return completion, nil
}


// chat without context
func GenerateContentLAI(api_token string, model_name string, promt string) (*llms.ContentResponse, error) {
	ctx := context.Background()
	token := api_token

	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(model_name),
		//llms.WithOptions()
		openai.WithBaseURL("http://localhost:8080/v1/"),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
	  log.Fatal(err)
	}
	

	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, "You are a helpfull assistant who help in whatever task human ask you about"),
		llms.TextParts(schema.ChatMessageTypeHuman, promt),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	return completion, nil
}



// chat with context without limitation of token to use
//  use it only to fast testing, REMOVE before production
func TestChatWithContextNoLimit(api_token string, model_name string) (string, error) {
	ctx := context.Background()
	token := api_token

	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(model_name),
		//llms.WithOptions()
		//openai.WithBaseURL("http://localhost:8080/v1/"),
		//openai.WithAPIVersion("v1"),
	)
	if err != nil {
	  log.Fatal(err)
	}

	memory_buffer := memory.NewConversationBuffer()

	//test data
	// First dialogue pair
	inputValues1 := map[string]any{"input": "Hi"}				// ignore linter
	outputValues1 := map[string]any{"output": "What's up"}

	memory_buffer.SaveContext(ctx,inputValues1,outputValues1)	//initial messages should be put like this

	memory_buffer.ChatHistory.AddUserMessage(ctx, "Not much, just hanging")  	// next messages from conversation could be added like this
	memory_buffer.ChatHistory.AddAIMessage(ctx,"Cool")
	memory_buffer.ChatHistory.AddUserMessage(ctx, "I am working at my new exiting golang AI project called 'Andromeda'")
	memory_buffer.ChatHistory.AddUserMessage(ctx, "My name is Bekket btw")
	
	conversation := chains.NewConversation(llm,memory_buffer) 	// build chain, start new conversation thread
	

	// Run is used when we have only one input (promt for example).   If there are need in passing few inputs then use chains.Call instead
	result, err := chains.Run(ctx,conversation,"what is my name and what project am I currently working on?")	//ignore linter error
	if err != nil {
		return "", err
	}

	// Example using call with few inputs
	/*
		translatePrompt := prompts.NewPromptTemplate(
		"Translate the following text from {{.inputLanguage}} to {{.outputLanguage}}. {{.text}}",
		[]string{"inputLanguage", "outputLanguage", "text"},
	)
	llmChain = chains.NewLLMChain(llm, translatePrompt)

	// Otherwise the call function must be used.
	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"inputLanguage":  "English",
		"outputLanguage": "French",
		"text":           "I love programming.",
	})
	if err != nil {
		return err
	}

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return fmt.Errorf("invalid chain return")
	}
	fmt.Println(out)
	*/

	log.Println("AI answer:")
	log.Println(result)

	log.Println("check if it's stored in messages, printing messages:")
	history, err := memory_buffer.ChatHistory.Messages(ctx)
	if err != nil {
		return "", err
	}
	//log.Println(history)
	total_turns := len(history)
	log.Println("total number of turns: ", total_turns)
	// Iterate over each message and print
    log.Println("Printing messages:")
    for _, msg := range history {
        log.Println(msg.GetContent())
    }

	return result,err
}


// Initialize New Dialog thread with User with no limitation for token usage (may fail, use with limit)
func InitializeNewChatWithContextNoLimit(api_token string, model_name string) (*ChatSession, error)  {
	//ctx := context.Background()

    llm, err := openai.New(
        openai.WithToken(api_token),
        openai.WithModel(model_name),
    )
    if err != nil {
        return nil, err
    }

    memoryBuffer := memory.NewConversationBuffer()
    conversation := chains.NewConversation(llm, memoryBuffer)

    return &ChatSession{
        ConversationBuffer: memoryBuffer,
        DialogThread: &conversation,
    }, nil
}


// Continue Dialog with memory included, so user can chat with remembering context of previouse messages
func ContinueChatWithContextNoLimit(session *ChatSession, prompt string) (string, error) {
	ctx := context.Background()
    result, err := chains.Run(ctx, session.DialogThread, prompt)
    if err != nil {
        return "", err
    }
    return result, nil
}



// TODO: remove or transfer this into tests
func TestOAI(api_token string)  {
	ctx := context.Background()
	token := api_token

	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel("gpt-3.5-turbo"),
	)
	if err != nil {
	  log.Fatal(err)
	}
	content := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(schema.ChatMessageTypeHuman, "What would be a good company name a company that makes colorful socks? Write at least 10 options"),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}



// TODO: make one function for both OAI & LAI, add baseUrl as argument
// Main function for generating from single promt (without memory and context)
func GenerateFromSinglePromtLocal(prompt string, model_name string) (string,error) {
	ctx := context.Background()
	llm, err := openai.New(
		//openai.WithToken()
		openai.WithBaseURL("http://localhost:8080"),
		openai.WithModel(model_name),
	)
	if err != nil {
	  log.Fatal(err)
	}
	
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
	 // log.Fatal(err)
	 return "", err
	}
	fmt.Println(completion)
	return completion, nil
}

func GenerateFromSinglePromtOAI(promt string, model_name string, api_token string) (string , error) {
	ctx := context.Background()
	llm, err := openai.New(
		openai.WithToken(api_token),
		openai.WithModel(model_name),
	)
	if err != nil {
	  log.Fatal(err)
	}
	
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, promt)
	if err != nil {
	 // log.Fatal(err)
	 return "", err
	}
	fmt.Println(completion)
	return completion, nil
}

