# trello-to-github-migrator

A simple tool to copy over all cards in a Trello board into a GitHub Project. The only metadata it copies over is the title of the trello card, which becomes the note text. 

## Prerequisites

- First off, create the Github repo and a project.
- Create all the lists you have in your Trello board in in your Github project

## Usage

- Run `glide install` to pull down dependencies
- Run `go build .` to run the tool
- Run the command with the following:

```
./trello-to-github-migrator -github-username <your-username> \
                            -github-access-token <get-at: https://github.com/settings/tokens> \
                            -github-owner <your-owner> \
                            -github-repo <your-repo> \
                            -github-project <your-project> \
                            -trello-username <your-username> \
                            -trello-app-key <get-at: https://trello.com/app-key> \
                            -trello-token <get-at: https://trello.com/app-key> \
                            -trello-board-name <your-board>
```
