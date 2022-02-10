package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func getPlayers() ([]string, error) {
	playersB, err := os.ReadFile("players.txt")
	if err != nil {
		return nil, err
	}
	players := strings.Split(strings.TrimSuffix(string(playersB), "\n"), "\n")
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	return players, nil
}

func getPrompt() (string, error) {
	promptID := rand.Int31n(49290) + 50000
	resp, err := http.Get(fmt.Sprintf("https://www.squibler.io/api/v1/random-prompt-generator/%d/", promptID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	prompt := new(struct {
		PromptContent string `json:"prompt_content"`
	})
	if err := json.Unmarshal(respBytes, &prompt); err != nil {
		return "", err
	}

	return prompt.PromptContent, nil
}

func randChar() string {
	return string(rune(rand.Int31n(26) + 65))
}

var types = []string{"noun", "verb", "adjective"}

func main() {
	rand.Seed(time.Now().Unix())

	prompt, err := getPrompt()
	if err != nil {
		panic(err)
	}
	fmt.Println("===== PROMPT =====")
	fmt.Println(prompt, "\n")

	players, err := getPlayers()
	if err != nil {
		panic(err)
	}

	smiths, agents := []string{}, []string{}
	for i, player := range players {
		if i%2 == 0 {
			smiths = append(smiths, player)
		} else {
			agents = append(agents, player)
		}
	}

	sort.Strings(smiths)
	sort.Strings(agents)

	fmt.Println("===== SMITHS =====")
	for i, smith := range smiths {
		fmt.Println(smith, ":", types[i%3], ":", randChar())
	}
	fmt.Println()

	fmt.Println("===== AGENTS =====")
	for _, agent := range agents {
		fmt.Println(agent)
	}
}
