package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	AriaLabel string `json:"ariaLabel"`
}

type CategorizedVideos struct {
	Category string  `json:"category"`
	Videos   []Video `json:"videos"`
}

func main() {
	// Read and parse scrape.json
	videos, err := readScrapeJSON("scrape.json")
	if err != nil {
		log.Fatalf("Error reading scrape.json: %v", err)
	}

	// Print the number of videos
	fmt.Printf("Number of videos: %d\n", len(videos))

	// Define category keywords
	categoryKeywords := map[string][]string{
		"Programming & Development":           {"Coded", "VS Code", "YAML", "programming", "development", "coding", "go", "python", "java", "javascript", "Devcontainers", "vscode", "visual studio code", "intellij", "eclipse", "netbeans", "atom", "sublime text", "vim", "emacs", "code editor", "ide", "integrated development environment", "developer", "Angular", "Node.js", "TypeScript", "Stripe"},
		"Cloud & Infrastructure":              {"cloud", "infrastructure", "aws", "azure", "gcp", "google", "cloud platform", "cloud services", "cloud computing", "cloud storage", "cloud networking", "cloud security", "cloud databases", "cloud migration", "cloud architecture", "cloud design", "cloud deployment", "cloud management", "cloud monitoring", "cloud scaling", "cloud optimization", "cloud performance", "cloud reliability", "cloud availability", "cloud fault tolerance", "cloud disaster recovery", "cloud backup", "cloud restore", "cloud pricing", "cloud billing", "cloud cost management", "cloud governance", "cloud compliance", "cloud audits", "cloud reviews", "cloud ratings", "cloud rankings", "cloud awards", "cloud recognition", "cloud certifications", "cloud badges", "cloud labels", "cloud tags", "cloud categories", "cloud topics", "cloud subjects", "cloud areas", "cloud domains", "cloud fields", "cloud industries", "cloud sectors", "cloud verticals", "cloud markets", "cloud audiences", "cloud users", "cloud developers", "cloud architects", "cloud engineers", "cloud administrators", "cloud operators", "cloud managers", "cloud directors", "cloud leads", "cloud officers", "cloud coordinators", "cloud specialists", "cloud consultants", "cloud advisors", "cloud partners", "cloud vendors", "cloud customers", "cloud clients", "cloud consumers", "cloud producers", "cloud providers", "cloud services", "cloud solutions", "cloud products", "cloud offerings", "cloud features", "cloud capabilities", "cloud integrations", "cloud extensions", "cloud plugins", "cloud modules", "cloud packages", "cloud dependencies", "cloud security", "cloud compliance", "cloud audits", "cloud reviews", "cloud ratings", "cloud rankings", "cloud awards", "cloud recognition", "cloud certifications", "cloud badges", "cloud labels", "cloud tags", "cloud categories", "cloud topics", "cloud subjects", "cloud areas", "cloud domains", "cloud fields", "cloud industries", "cloud sectors", "cloud verticals", "cloud markets", "cloud audiences", "cloud users", "cloud developers", "cloud architects", "cloud engineers", "cloud administrators", "cloud operators", "cloud managers", "cloud directors", "cloud leads", "cloud officers", "cloud coordinators", "cloud specialists", "cloud consultants", "cloud advisors", "cloud partners", "cloud vendors", "cloud customers", "cloud clients", "cloud consumers"},
		"DevOps and CI/CD":                    {"Vault", "Ansible", "secrets management", "HashiCorp", "Secrets Management", "devops", "ci/cd", "continuous integration", "continuous delivery", "terraform", "platform engineering", "site reliability engineering", "sre", "devops engineer", "devops architect", "devops consultant", "devops specialist", "devops manager", "devops director", "devops lead", "devops officer", "devops coordinator", "devops tools", "devops practices", "devops principles", "devops culture", "devops automation", "devops monitoring", "devops scaling", "devops optimization", "devops performance", "devops reliability", "devops availability", "devops fault tolerance", "devops disaster recovery", "devops backup", "devops restore", "devops security", "devops compliance", "devops audits", "devops reviews", "devops ratings", "devops rankings", "devops awards", "devops recognition", "devops certifications", "devops badges", "devops labels", "devops tags", "devops categories", "devops topics", "devops subjects", "devops areas", "devops domains", "devops fields", "devops industries", "devops sectors", "devops verticals", "devops markets", "devops audiences", "devops users", "devops developers", "devops architects", "devops engineers", "devops administrators", "devops operators", "devops managers", "devops directors", "devops leads", "devops officers", "devops coordinators", "devops specialists", "devops consultants", "devops advisors", "devops partners", "devops vendors", "devops customers", "devops clients", "devops consumers", "devops producers", "devops providers", "devops services", "devops solutions", "devops products", "devops offerings", "devops features", "devops capabilities", "devops integrations", "devops extensions", "devops plugins", "devops modules", "devops packages", "devops dependencies", "devops security", "devops compliance", "devops audits", "devops reviews", "devops ratings", "devops rankings", "devops awards", "devops recognition", "devops certifications", "devops badges", "devops labels", "devops tags", "devops categories", "devops topics", "devops subjects", "devops areas", "devops domains"},
		"Containers and Kubernetes":           {"Operators", "Talos", "talos", "KubeCon", "Stateful", "microservices", "Helm", "Knative", "OpenShift", "Open Policy Agent", "K8s", "containers", "kubernetes", "docker", "Kubernetes", "Docker", "containerization", "container orchestration", "container management", "container deployment", "container scaling", "container optimization", "container performance", "container reliability", "container availability", "container fault tolerance", "container disaster recovery", "container backup", "container restore", "container security", "container compliance", "container audits", "container reviews", "container ratings", "container rankings", "container awards", "container recognition", "container certifications", "container badges", "container labels", "container tags", "container categories", "container topics", "container subjects", "container areas", "container domains", "container fields", "container industries", "container sectors", "container verticals", "container markets", "container audiences", "container users", "container developers", "container architects", "container engineers", "container administrators", "container operators", "container managers", "container directors", "container leads", "container officers", "container coordinators", "container specialists", "container consultants", "container advisors", "container partners", "container vendors", "container customers", "container clients", "container consumers", "container producers", "container providers", "container services", "container solutions", "container products", "container offerings", "container features", "container capabilities", "container integrations", "container extensions", "container plugins", "container modules", "container packages", "container dependencies", "container security", "container compliance", "container audits", "container reviews", "container ratings", "container rankings", "container awards", "container recognition", "container certifications", "container badges", "container labels", "container tags", "container categories", "container topics", "container subjects", "container areas", "container domains", "container fields", "container industries", "container sectors", "container verticals", "container markets", "container audiences", "container users", "container developers", "container architects", "container engineers", "container administrators", "container operators", "container managers", "container directors", "container leads", "container officers", "container coordinators", "container specialists", "container consultants", "container advisors", "container partners", "container vendors", "container customers", "container clients", "container consumers", "container producers", "container providers", "container services", "container solutions", "container products", "container offerings", "container features", "container capabilities", "container integrations", "container extensions", "container plugins", "container modules", "container packages"},
		"Data Management and Databases":       {"schema", "Schemas", "DB", "Data Protection", "SurrealDB", "Disaster Recovery", "Storage", "data management", "databases", "sql", "nosql", "mongodb", "postgresql", "MongoDB", "PostgreSQL", "MySQL", "MariaDB", "Cassandra", "Couchbase", "CouchDB", "DynamoDB", "Aurora", "RDS", "Redshift", "BigQuery", "Snowflake", "database"},
		"Cloud-Native and Serverless":         {"Fermyon", "Service Mesh", "cloud-native", "serverless", "lambda", "functions", "cloud functions", "serverless framework", "cloud run", "cloudflare", "faas", "paas", "saas", "iaas", "cloud computing"},
		"Security and DevSecOps":              {"Hack", "NIS2", "security", "devsecops", "cybersecurity", "infosec", "information security", "security engineering", "security operations", "security architecture", "security analyst", "security consultant", "security specialist", "security engineer", "security architect", "security operations center", "security operations centre", "security operations analyst", "security operations engineer", "security operations architect", "security operations specialist", "security operations consultant", "security operations manager", "security operations director", "security operations lead", "security operations officer", "security operations coordinator"},
		"Open Source and Community":           {"Open-Source", "open source", "community", "opensource", "github", "gitlab", "bitbucket", "source control", "version control", "git", "versioning", "gitops", "gitflow", "github actions", "gitlab ci", "bitbucket pipelines", "open source software", "open source projects", "open source contributions", "open source development", "open source licensing", "open source governance", "open source community", "open source ecosystem", "open source tools", "open source technologies", "open source frameworks", "open source libraries", "open source modules", "open source packages", "open source dependencies", "open source security", "open source compliance", "open source audits", "open source reviews", "open source releases", "open source updates", "open source patches", "open source bug fixes", "open source enhancements", "open source features", "open source requests", "open source contributions", "open source pull requests", "open source issues", "open source discussions", "open source forums", "open source chats", "open source meetups", "open source events", "open source conferences", "open source summits", "open source workshops", "open source tutorials", "open source webinars", "open source videos", "open source podcasts", "open source blogs", "open source articles", "open source books", "open source papers", "open source research", "open source studies", "open source surveys", "open source polls", "open source feedback", "open source reviews", "open source ratings", "open source rankings", "open source awards", "open source recognition", "open source certifications", "open source badges", "open source labels", "open source tags", "open source categories", "open source topics", "open source subjects", "open source areas", "open source domains", "open source fields", "open source industries", "open source sectors", "open source verticals", "open source markets", "open source audiences", "open source users", "open source developers", "open source contributors", "open source maintainers", "open source reviewers", "open source approvers", "open source committers", "open source authors", "open source editors", "open source publishers", "open source consumers", "open source producers", "open source providers", "open source consumers", "open source customers", "open source clients", "open source partners", "open source vendors"},
		"Storytelling and Career Development": {"CTO", "CEO", "Story", "story", "Job", "storytelling", "story telling", "career development", "career", "development", "career growth", "career advancement", "career progression", "career success", "career satisfaction", "career fulfillment", "career happiness", "career wellbeing", "career balance", "career stability", "career security", "career safety", "career health", "career wealth", "career prosperity", "career abundance", "career opportunities", "career options", "career choices", "career decisions", "career planning", "career strategy", "career management", "career leadership", "career mentorship", "career coaching", "career training", "career education", "career learning", "career development", "career growth", "career advancement", "career progression", "career success", "career satisfaction", "career fulfillment", "career happiness", "career wellbeing", "career balance", "career stability", "career security", "career safety", "career health", "career wealth", "career prosperity", "career abundance", "career opportunities", "career options", "career choices", "career decisions", "career planning", "career strategy", "career management", "career leadership", "career mentorship", "career coaching", "career training", "career education", "career learning", "career development", "career growth", "career advancement", "career progression", "career success", "career satisfaction", "career fulfillment", "career happiness", "career wellbeing", "career balance", "career stability", "career security", "career safety", "career health", "career wealth", "career prosperity", "career abundance", "career opportunities", "career options", "career choices", "career decisions", "career planning", "career strategy", "career management", "career leadership", "career mentorship", "career coaching", "career training", "career education", "career learning", "career development", "career growth", "career advancement", "career progression", "career success", "career satisfaction", "career fulfillment", "career happiness", "career wellbeing", "career balance", "career stability", "career security", "career safety", "career health", "career wealth", "career prosperity", "career abundance", "career opportunities", "career options", "career choices", "career decisions", "career planning", "career strategy", "career management", "career leadership", "career mentorship", "career coaching", "career training", "career education", "career learning", "career development", "career growth", "career advancement", "career progression"},
		"AI and Emerging Technologies":        {"Ai", "AI", "GPT", "LLM", "Ollama", "GPT-3", "ai", "artificial intelligence", "machine learning", "emerging technologies", "blockchain", "quantum computing", "iot", "internet of things", "edge computing", "fog computing", "distributed computing", "distributed systems", "distributed systems design", "distributed systems architecture", "distributed systems engineering", "distributed systems development", "distributed systems operations", "distributed systems management", "distributed systems monitoring", "distributed systems testing", "distributed systems deployment", "distributed systems scaling", "distributed systems performance", "distributed systems optimization", "distributed systems security", "distributed systems reliability", "distributed systems availability", "distributed systems fault tolerance", "distributed systems disaster recovery", "distributed systems backup", "distributed systems restore"},
		"Tools and Productivity":              {"Tmux", "Canva", "tools", "productivity", "efficiency", "automation", "tooling", "toolchain", "toolset", "toolkit", "toolbox", "toolbelt", "toolbox", "toolbelt", "toolkit", "toolchain", "toolset", "tooling", "automation", "efficiency", "productivity", "tools and productivity", "productivity and tools", "tools for productivity", "productivity tools", "tools for efficiency", "efficiency tools", "tools for automation", "automation tools", "tools for tooling", "tooling tools", "tools for toolchain", "toolchain tools", "tools for toolset", "toolset tools", "tools for toolkit", "toolkit tools", "tools for toolbox", "toolbox tools", "tools for toolbelt", "toolbelt tools", "tools for toolbox", "toolbox tools", "tools for toolbelt", "toolbelt tools", "tools for productivity and efficiency", "productivity and efficiency tools", "tools for productivity and automation", "productivity and automation tools", "tools for productivity and tooling", "productivity and tooling tools", "tools for productivity and toolchain", "productivity and toolchain tools", "tools for productivity and toolset", "productivity and toolset tools", "tools for productivity and toolkit", "productivity and toolkit tools", "tools for productivity and toolbox", "productivity and toolbox tools", "tools for productivity and toolbelt", "productivity and toolbelt tools", "tools for efficiency and automation", "efficiency and automation tools", "tools for efficiency and tooling", "efficiency and tooling tools", "tools for efficiency and toolchain", "efficiency and toolchain tools", "tools for efficiency and toolset", "efficiency and toolset tools", "tools for efficiency and toolkit", "efficiency and toolkit tools", "tools for efficiency and toolbox", "efficiency and toolbox tools", "tools for efficiency and toolbelt", "efficiency and toolbelt tools", "tools for automation and tooling", "automation and tooling tools", "tools for automation and toolchain", "automation and toolchain tools", "tools for automation and toolset", "automation and toolset tools", "tools for automation and toolkit", "automation and toolkit tools", "tools for automation and toolbox", "automation and toolbox tools", "tools for automation and toolbelt", "automation and toolbelt tools", "tools for tool"},
		"Linux":                               {"linux", "Linux", "ubuntu", "debian", "centos", "redhat", "fedora", "suse", "arch", "manjaro", "mint", "elementary", "popos", "kali", "raspbian", "raspberrypi", "raspberry pi", "raspberry", "pi", "linux kernel", "linux distributions", "linux distros", "linux desktop", "linux server", "linux laptop", "linux workstation", "linux desktop environment", "linux window manager", "linux shell", "linux terminal", "linux command line", "linux bash", "linux zsh", "linux fish", "linux ksh", "linux csh", "linux tcsh", "linux sh", "linux scripting", "linux programming", "linux development", "linux administration", "linux operations", "linux management", "linux monitoring", "linux scaling", "linux optimization", "linux performance", "linux reliability", "linux availability", "linux fault tolerance", "linux disaster recovery", "linux backup", "linux restore", "linux security", "linux compliance", "linux audits", "linux reviews", "linux ratings", "linux rankings", "linux awards", "linux recognition", "linux certifications", "linux badges", "linux labels", "linux tags", "linux categories", "linux topics", "linux subjects", "linux areas", "linux domains", "linux fields", "linux industries", "linux sectors", "linux verticals", "linux markets", "linux audiences", "linux users", "linux developers", "linux architects", "linux engineers", "linux administrators", "linux operators", "linux managers", "linux directors", "linux leads", "linux officers", "linux coordinators", "linux specialists", "linux consultants", "linux advisors", "linux partners", "linux vendors", "linux customers", "linux clients", "linux consumers", "linux producers", "linux providers", "linux services", "linux solutions", "linux products", "linux offerings", "linux features", "linux capabilities", "linux integrations", "linux extensions", "linux plugins", "linux modules", "linux packages", "linux dependencies", "linux security", "linux compliance", "linux audits", "linux reviews", "linux ratings", "linux rankings", "linux awards", "linux recognition", "linux certifications", "linux badges", "linux labels", "linux tags", "linux categories"},
		"Virtualisation":                      {"virtualisation", "virtualization", "vm", "vmware", "virtualbox", "hypervisor", "kvm", "xen", "qemu", "VMware", "Proxmox", "proxmox", "esxi", "vcenter", "vSphere", "vSAN", "vRealize", "vCloud", "vCloud Director", "vCloud Suite", "vCloud Air", "vCloud Hybrid Service", "vCloud Connector", "vCloud Networking and Security", "vCloud Automation Center", "vCloud Application Director", "vCloud Operations Management Suite", "vCloud Suite SDK", "Hyperv", "hyperv", "hyper-v"},
	}

	// Categorize videos
	categorizedVideos := categorizeVideos(videos, categoryKeywords)

	// Save categorized videos to a new JSON file
	if err := saveCategorizedVideos("categorized_videos.json", categorizedVideos); err != nil {
		log.Fatalf("Error saving categorized videos: %v", err)
	}

	fmt.Println("Categorized videos saved to categorized_videos.json")

	// Authenticate with YouTube Data API
	client, err := getClient("credentials.json")
	if err != nil {
		log.Fatalf("Error getting YouTube client: %v", err)
	}

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	// Create playlists for each category and add videos
	for category := range categoryKeywords {
		if err := createYouTubePlaylist(service, category, categorizedVideos); err != nil {
			log.Fatalf("Error creating YouTube playlist for category %s: %v", category, err)
		}
	}

	fmt.Println("YouTube playlists created for each category")
}

// readScrapeJSON reads and parses the scrape.json file
func readScrapeJSON(filename string) ([]Video, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var videos []Video
	if err := json.Unmarshal(bytes, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

// categorizeVideos categorizes videos based on their titles using a map of category keywords
func categorizeVideos(videos []Video, categoryKeywords map[string][]string) []CategorizedVideos {
	// Define the order of categories based on their specificity
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
		"Tools and Productivity",
		"Linux",
		"Virtualisation",
	}

	categorized := make([]CategorizedVideos, len(categories)+1)
	for i, category := range categories {
		categorized[i] = CategorizedVideos{Category: category}
	}
	categorized[len(categories)] = CategorizedVideos{Category: "Other"}

	for _, video := range videos {
		found := false
		for i, category := range categories {
			for _, keyword := range categoryKeywords[category] {
				if strings.Contains(strings.ToLower(video.Title), strings.ToLower(keyword)) {
					categorized[i].Videos = append(categorized[i].Videos, video)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			categorized[len(categories)].Videos = append(categorized[len(categories)].Videos, video)
		}
	}

	// Debug: Print categorized videos
	for _, catVideos := range categorized {
		fmt.Printf("Category: %s, Number of Videos: %d\n", catVideos.Category, len(catVideos.Videos))
	}

	return categorized
}

// saveCategorizedVideos saves the categorized videos to a JSON file
func saveCategorizedVideos(filename string, categorizedVideos []CategorizedVideos) error {
	bytes, err := json.MarshalIndent(categorizedVideos, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
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

// createYouTubePlaylist creates a playlist for a given category and adds videos to it
func createYouTubePlaylist(service *youtube.Service, category string, categorizedVideos []CategorizedVideos) error {
	for _, catVideos := range categorizedVideos {
		if catVideos.Category == category {
			// Create playlist
			playlist := &youtube.Playlist{
				Snippet: &youtube.PlaylistSnippet{
					Title:       category + " Playlist",
					Description: "A playlist of " + category + " videos",
				},
				Status: &youtube.PlaylistStatus{
					PrivacyStatus: "private",
				},
			}

			playlistResponse, err := service.Playlists.Insert([]string{"snippet", "status"}, playlist).Do()
			if err != nil {
				return fmt.Errorf("error creating playlist: %v", err)
			}

			// Add videos to playlist
			for _, video := range catVideos.Videos {
				videoID := extractVideoID(video.Link)
				fmt.Printf("Adding video to playlist: %s (ID: %s)\n", video.Title, videoID)

				playlistItem := &youtube.PlaylistItem{
					Snippet: &youtube.PlaylistItemSnippet{
						PlaylistId: playlistResponse.Id,
						ResourceId: &youtube.ResourceId{
							Kind:    "youtube#video",
							VideoId: videoID,
						},
					},
				}

				_, err := service.PlaylistItems.Insert([]string{"snippet"}, playlistItem).Do()
				if err != nil {
					return fmt.Errorf("error adding video to playlist: %v", err)
				}
			}

			fmt.Printf("Playlist created for category %s: %s\n", category, playlistResponse.Id)
			break
		}
	}
	return nil
}

// extractVideoID extracts the video ID from a YouTube link
func extractVideoID(link string) string {
	parts := strings.Split(link, "v=")
	if len(parts) > 1 {
		videoID := strings.Split(parts[1], "&")[0]
		fmt.Printf("Extracted video ID: %s from link: %s\n", videoID, link)
		return videoID
	}
	fmt.Printf("Failed to extract video ID from link: %s\n", link)
	return ""
}
