package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

const webhookURL = "https://discord.com/api/webhooks/1327300726657126420/sAVKp77UNSOvFp-ld5414aYBPmihyqg7dER8RvWUW2ZHtc0ILgmDBFmpUeCbo2kTwwe8" // Replace with your webhook URL

// cryptoRandIntn generates a random integer in the range [0, n) using crypto/rand
func cryptoRandIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid range: n must be greater than 0")
	}
	max := big.NewInt(int64(n))
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(randomInt.Int64()), nil
}

func generatePassword(length int, uppercaseCount int, specialcCount int, numberCount int) string {
	const lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	const uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?/"
	const numberChars = "0123456789"

	var password strings.Builder

	// Add uppercase characters
	for i := 0; i < uppercaseCount; i++ {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			log.Fatalf("Error generating random byte: %v", err)
		}
		password.WriteByte(uppercaseChars[randomByte[0]%byte(len(uppercaseChars))])
	}

	// Add special characters
	for i := 0; i < specialcCount; i++ {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			log.Fatalf("Error generating random byte: %v", err)
		}
		password.WriteByte(specialChars[randomByte[0]%byte(len(specialChars))])
	}

	// Add numbers
	for i := 0; i < numberCount; i++ {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			log.Fatalf("Error generating random byte: %v", err)
		}
		password.WriteByte(numberChars[randomByte[0]%byte(len(numberChars))])
	}

	// Fill the rest of the password length with lowercase characters
	for i := password.Len(); i < length; i++ {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			log.Fatalf("Error generating random byte: %v", err)
		}
		password.WriteByte(lowercaseChars[randomByte[0]%byte(len(lowercaseChars))])
	}

	// Shuffle the password to ensure randomness
	shuffledPassword := []rune(password.String())
	for i := len(shuffledPassword) - 1; i > 0; i-- {
		j, err := cryptoRandIntn(i + 1)
		if err != nil {
			log.Fatalf("Error generating random index: %v", err)
		}
		shuffledPassword[i], shuffledPassword[j] = shuffledPassword[j], shuffledPassword[i]
	}

	return string(shuffledPassword)
}

func sendToDiscord(password string) {
	message := map[string]string{"content": password}
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	_, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending message to Discord: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		length, err := strconv.Atoi(r.FormValue("length"))
		if err != nil {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			return
		}

		uppercaseCount := 0
		if r.FormValue("uppercaseCheckbox") == "on" {
			uppercaseCount, err = strconv.Atoi(r.FormValue("uppercaseCount"))
			if err != nil {
				http.Error(w, "Invalid uppercase count", http.StatusBadRequest)
				return
			}
		}

		specialcCount := 0
		if r.FormValue("specialcCheckbox") == "on" {
			specialcCount, err = strconv.Atoi(r.FormValue("specialcCount"))
			if err != nil {
				http.Error(w, "Invalid special character count", http.StatusBadRequest)
				return
			}
		}

		numberCount := 0
		if r.FormValue("numberCheckbox") == "on" {
			numberCount, err = strconv.Atoi(r.FormValue("numberCount"))
			if err != nil {
				http.Error(w, "Invalid number count", http.StatusBadRequest)
				return
			}
		}

		password := generatePassword(length, uppercaseCount, specialcCount, numberCount)
		sendToDiscord(password)

		// Redirect to the root URL to refresh the page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Handle password display
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query().Get("password")
		if password == "" {
			http.Error(w, "No password generated", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "<h1>Generated Password</h1><p>%s</p>", password)
	})

	http.ListenAndServe(":12500", nil)
}
