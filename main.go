package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"context"
	"strings"
	"bufio"

	"google.golang.org/genai"
)

func main() {
	f := flag.String("f", "", "Filename to save output")
	prompt := flag.String("p", "", "Prompt")
	key := flag.String("key", "", "API key")

	flag.Parse()

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Could not get user config dir:", err)
		os.Exit(1)
	}

	appDir := filepath.Join(configDir, "termini")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		fmt.Println("Could not create config dir:", err)
		os.Exit(1)
	}

	keyFile := filepath.Join(appDir, "key.txt")

	if *key != "" {
		err := os.WriteFile(keyFile, []byte(*key), 0600)
		if err != nil {
			fmt.Println("Error writing key file:", err)
			os.Exit(1)
		}
		fmt.Println("API key saved to", keyFile)
		return
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		fmt.Println("key.txt not found. Please provide a key with --key")
		os.Exit(1)
	}

	data, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error reading key file:", err)
		os.Exit(1)
	}
	apiKey := string(data)

	if *prompt == "" {
		fmt.Println("Please provide a prompt with --p")
		return
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	finalPrompt := ""

	if *f == "" {
		finalPrompt = "Respond with only the Linux command, no explanation, and do NOT use backticks or any code formatting:\n" + *prompt
	} else {
		finalPrompt = "Generate the full content of a file based on this description, and do NOT use backticks or any code formatting:\n" + *prompt
	}

	resp, err := client.Models.GenerateContent(context.Background(), "models/gemini-2.5-flash", []*genai.Content{
		{Parts: []*genai.Part{
			{Text: finalPrompt},
		}},
	}, nil)
	if err != nil {
		fmt.Println("Error generating content:", err)
		return
	}

	if resp.Text() != "" {
		if *f == "" {
			fmt.Println("Generated command:")
			fmt.Println(resp.Text())
			fmt.Println("Run it? Y/n")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			input = strings.TrimSpace(input)

			if input == "Y" {
				//cmdStr := strings.Trim(resp.Text(), "`")

				parts := strings.Fields(resp.Text())
				if len(parts) == 0 {
					fmt.Println("No command to run.")
					return
				}

				cmd := exec.Command(parts[0], parts[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					fmt.Println("Error running command:", err)
				}
			} else {
				return
			}
		} else {
			//data := strings.Trim(resp.Text(), "`")
			err := os.WriteFile(*f, []byte(resp.Text()), 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	} else {
		fmt.Println("No response from Gemini")
	}
}
