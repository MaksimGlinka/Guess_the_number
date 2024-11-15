package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GameData struct {
	Wins        int
	Attempts    int
	Message     string
	GameStarted bool
	Level       int
}

var (
	targetNumber int
	gameData     GameData
	templates    *template.Template
)

func init() {
	rand.Seed(time.Now().UnixNano())
	funcMap := template.FuncMap{
		"contains": strings.Contains,
	}
	templates = template.Must(template.New("game.html").Funcs(funcMap).ParseFiles("game.html"))

}

func main() {
	http.HandleFunc("/", handleGame)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/guess", handleGuess)

	http.ListenAndServe(":8080", nil)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "game.html", gameData)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	level, _ := strconv.Atoi(r.FormValue("level"))
	gameData.Level = level
	gameData.GameStarted = true
	gameData.Attempts = 3
	gameData.Message = "игра началась у вас три попытки"

	switch level {
	case 1:
		targetNumber = rand.Intn(11)
	case 2:
		targetNumber = rand.Intn(51)
	case 3:
		targetNumber = rand.Intn(101)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleGuess(w http.ResponseWriter, r *http.Request) {
	if !gameData.GameStarted {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	guess, _ := strconv.Atoi(r.FormValue("guess"))
	gameData.Attempts--

	if guess == targetNumber {
		gameData.Wins++
		gameData.Message = "Поздравляем! Вы угадали число!"
		gameData.GameStarted = false
	} else if gameData.Attempts == 0 {
		gameData.Message = "Игра окончена! Загаданное число было: " + strconv.Itoa(targetNumber)
		gameData.GameStarted = false
	} else if guess < targetNumber {
		gameData.Message = "Загаданное число больше! Осталось попыток: " + strconv.Itoa(gameData.Attempts)
	} else {
		gameData.Message = "Загаданное число меньше! Осталось попыток: " + strconv.Itoa(gameData.Attempts)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
