package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/x0ptr/aidoc/format"
	"github.com/x0ptr/aidoc/storage"
)

func Run() error {
	if CLI.Verbose {
		fmt.Printf("[verbose] language=%s topic=%s\n", CLI.Language, CLI.Topic)
	}
	if !CLI.SkipCache {
		if ans, ok := storage.CacheGet(CLI.Language, CLI.Topic, CLI.Verbose); ok {
			mustNil(print(ans, CLI.Output, true))
			return nil
		}
	}

	answer := concatStrings(completion(CLI.Language, CLI.Topic, CLI.Verbose), CLI.Language, CLI.Topic, CLI.Output)
	storage.CacheSet(CLI.Language, CLI.Topic, answer, CLI.Verbose)
	mustNil(storage.SaveCache())
	print(answer, CLI.Output, false)
	return nil
}

var CLI struct {
	Verbose   bool   `short:"v" help:"Enable verbose output."`
	Output    bool   `short:"o" help:"Write output directly to stdout instead of pager."`
	SkipCache bool   `short:"s" help:"Skip cache."`
	Language  string `arg:"" name:"language" help:"the programming language"`
	Topic     string `arg:"" name:"topic"    help:"the topic you want to lookup"`
}

func completion(language, topic string, isVerbose bool) string {
	verbose := ""
	if isVerbose {
		verbose = "VERBOSE"
	}
	cfg, err := storage.LoadConfig()
	client := openai.NewClient(
		option.WithAPIKey(cfg.OpenAIKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a very specific documentation machine. " +
				"You can't talk ONLY give information about the programming language and topic i provide you." +
				"When I don't send you 'VERBOSE' you can't be verbose. Answer realy short answers. A few lines of explenation and one short example."),
			openai.UserMessage(fmt.Sprintf("Programming Language: %s\nTopic: %s -- %s\n", language, topic, verbose)),
		},
		Model: openai.ChatModelGPT5ChatLatest,
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}

func concatStrings(answer, language, topic string, skipColor bool) string {
	if !skipColor {
		return answer
	}
	result := format.PrintBanner(fmt.Sprintf("%s - %s", language, topic))
	result = fmt.Sprintf("%s\n%s", result, answer)
	return result
}

func print(result string, output, fromCache bool) error {
	if fromCache {
		result += "\n-- FROM CACHE --"
	}
	if !output {
		pager := os.Getenv("PAGER")
		if pager == "" {
			pager = "less"
		}

		cmd := exec.Command(pager)
		cmd.Stdin = io.NopCloser(strings.NewReader(result))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
	println(result)
	return nil
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--clear-cache" {
		storage.ClearCache()
		log.Println("cleared cache.")
		os.Exit(0)
	}

	if len(os.Args) == 3 && os.Args[1] == "--set-apikey" {
		storage.SaveAPIKey(os.Args[2])
		log.Println("api key stored.")
		os.Exit(0)
	}

	must(storage.CacheFilePath())
	must(storage.LoadCache())
	must(storage.ConfigFilePath())

	ctx := kong.Parse(&CLI,
		kong.Name("aidoc"),
		kong.Description("AI backed doc helper.\nUse --clear-cache to wipe the cache\n--set-apikey [key] for saving openai api key"),
		kong.UsageOnError(),
	)

	err := Run()
	ctx.FatalIfErrorf(err)
	os.Exit(0)
}
