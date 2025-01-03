package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	// Authenticate with YouTube Data API
	client, err := getClient("credentials.json")
	if err != nil {
		log.Fatalf("Error getting YouTube client: %v", err)
	}

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	// Define categories
	categories := []string{
		"Programming & Development",
		"Cloud & Infrastructure",
		"DevOps and CI/CD",
		"Containers and Kubernetes",
		"Data Management and Databases",
		"Cloud-Native and Serverless",
		"Security and DevSecOps",
		"Open Source and Community",
		"Storytelling and Career Development",
		"AI and Emerging Technologies",
		"Workshops and Tutorials",
		"Tools and Productivity",
	}

	// Delete playlists for each category
	for {
		deleted, err := deletePlaylists(service, categories)
		if err != nil {
			log.Fatalf("Error deleting playlists: %v", err)
		}
		if !deleted {
			break
		}
		fmt.Println("Waiting for quota reset...")
		time.Sleep(1 * time.Minute) // Wait for 1 minute before the next batch
	}

	fmt.Println("Playlists deleted successfully")
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(credentialsFile string) (*http.Client, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	token, err := tokenFromFile("token.json")
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken("token.json", token)
	}

	return config.Client(context.Background(), token), nil
}

// tokenFromFile retrieves a Token from a given file path.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// deletePlaylists lists and deletes playlists that match the given categories
func deletePlaylists(service *youtube.Service, categories []string) (bool, error) {
	call := service.Playlists.List([]string{"id", "snippet"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return false, fmt.Errorf("error listing playlists: %v", err)
	}

	deleted := false
	for _, playlist := range response.Items {
		for _, category := range categories {
			if strings.Contains(playlist.Snippet.Title, category) {
				fmt.Printf("Deleting playlist: %s (ID: %s)\n", playlist.Snippet.Title, playlist.Id)
				if err := service.Playlists.Delete(playlist.Id).Do(); err != nil {
					return false, fmt.Errorf("error deleting playlist: %v", err)
				}
				deleted = true
				break
			}
		}
	}

	return deleted, nil
}
