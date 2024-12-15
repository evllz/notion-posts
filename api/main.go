package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// func Main() {
// 	fmt.Println("Notion Posts API Ver. 1.0.0")
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "API Server listening on port 8080")
// 	})

// 	mux.HandleFunc("GET /posts", getPosts)

// 	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

func GetPosts(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Validate the token (example assumes Bearer token format)
	authToken := os.Getenv("AUTH_TOKEN") // Replace with your secret or validation logic
	expectedToken := fmt.Sprintf("Bearer %s", authToken)
	if authHeader != expectedToken {
		http.Error(w, "Invalid or unauthorized token", http.StatusUnauthorized)
		return
	}

	//Notion API Details

	notionDatabaseId := os.Getenv("NOTION_DATABASE_ID")
	notionToken := os.Getenv("NOTION_TOKEN")
	notionApiUrl := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", notionDatabaseId)

	// Prepare the request body (empty in this case)
	requestBody, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		http.Error(w, "Failed to prepare request body", http.StatusInternalServerError)
		return
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", notionApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to create HTTP request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send HTTP request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check for HTTP status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Non-OK HTTP status: %d, body: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
		return
	}

	// Parse the response body
	var responseData struct {
		Results []interface{} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	// Convert the results to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData.Results)
}
