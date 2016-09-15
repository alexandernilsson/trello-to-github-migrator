package main

import (
	"flag"
	"log"

	trello "github.com/VojtechVitek/go-trello"
)

var (
	githubUsername    = flag.String("github-username", "", "Github username")
	githubAccessToken = flag.String("github-access-token", "", "Github access token")
	githubOwner       = flag.String("github-owner", "", "Github organization owner/name")
	githubRepo        = flag.String("github-repo", "", "Repo to create cards in")
	githubProject     = flag.String("github-project", "", "Github project to copy cards to")

	trelloUsername  = flag.String("trello-username", "", "Trello username")
	trelloAppKey    = flag.String("trello-app-key", "", "Trello application key")
	trelloToken     = flag.String("trello-token", "", "Trello application token")
	trelloBoardName = flag.String("trello-board-name", "", "Trello board to migrate")
)

func init() {
	flag.Parse()
}

func main() {
	// create new Github client
	github := GithubProjectAPI{
		Username:    *githubUsername,
		AccessToken: *githubAccessToken,
		Owner:       *githubOwner,
		Repo:        *githubRepo,
	}

	// get the Github project
	project, err := github.GetGitHubProject(*githubProject)
	if err != nil {
		log.Fatal(err)
	}

	// create new Trello Client
	trello, err := trello.NewAuthClient(*trelloAppKey, trelloToken)
	if err != nil {
		log.Fatal(err)
	}

	// fetch user info
	log.Printf("fetching trello user %s", *trelloUsername)
	user, err := trello.Member(*trelloUsername)
	if err != nil {
		log.Fatal(err)
	}

	// for users boards
	log.Printf("fetching trello boards for user %s", *trelloUsername)
	boards, err := user.Boards()
	if err != nil {
		log.Fatal(err)
	}

	// loop over all boards, find the matching one
	for _, board := range boards {
		if board.Name == *trelloBoardName {
			log.Printf("found matching board %s", *trelloBoardName)

			// fetch all lists for that board
			log.Printf("fetching lists for board %s", *trelloBoardName)
			lists, err := board.Lists()
			if err != nil {
				log.Fatal(err)
			}

			// iterate over lists in trello board
			for _, list := range lists {
				// find matching column in the github project
				column, err := github.GetGitHubProjectColumn(project.Number, list.Name)
				if err != nil {
					log.Fatal(err)
				}

				// list all cards in list
				cards, _ := list.Cards()
				for _, card := range cards {
					log.Printf("creating %s in github column %s\n", card.Name, column.Name)

					github.CreateGithubProjectCard(column, card.Name)
				}
			}
		}
	}
}
