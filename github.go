package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GitHubProject holds the relevant stuff from getting a project via API
type GitHubProject struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

// GitHubProjectColumn holds the relevant stuff from getting a project column via API
type GitHubProjectColumn struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateGitHubProjectCard holds all relevant stuff to build a request to create a card
type CreateGitHubProjectCard struct {
	Note        string `json:"note,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

// GithubProjectAPI holds all the goodies that's needed to interact with the API
type GithubProjectAPI struct {
	Username    string
	AccessToken string
	Owner       string
	Repo        string
}

// wrapper around a http client, doing some github api specific stuff
func (g *GithubProjectAPI) githubRequest(method, path, body string) ([]byte, error) {

	url := fmt.Sprintf("https://%s:%s@api.github.com%s", g.Username, g.AccessToken, path)
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.inertia-preview+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetGitHubProject gets a github project by it's name
func (g *GithubProjectAPI) GetGitHubProject(name string) (GitHubProject, error) {
	resp, err := g.githubRequest("GET", fmt.Sprintf("/repos/%s/%s/projects", g.Owner, g.Repo), "")
	if err != nil {
		return GitHubProject{}, err
	}

	var projects []GitHubProject
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return GitHubProject{}, err
	}
	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}
	return GitHubProject{}, fmt.Errorf("Project %s not found", name)
}

// GetGitHubProjectColumn gets a github project column
func (g *GithubProjectAPI) GetGitHubProjectColumn(number int, name string) (GitHubProjectColumn, error) {
	resp, err := g.githubRequest("GET", fmt.Sprintf("/repos/%s/%s/projects/%d/columns", g.Owner, g.Repo, number), "")
	if err != nil {
		return GitHubProjectColumn{}, err
	}

	var columns []GitHubProjectColumn
	err = json.Unmarshal(resp, &columns)
	if err != nil {
		return GitHubProjectColumn{}, err
	}
	for _, column := range columns {
		if column.Name == name {
			return column, nil
		}
	}
	return GitHubProjectColumn{}, fmt.Errorf("Column %s not found", name)
}

// CreateGithubProjectCard creates a github project card in a column
func (g *GithubProjectAPI) CreateGithubProjectCard(column GitHubProjectColumn, name string) {
	projectCard := CreateGitHubProjectCard{
		Note:        name,
		ContentType: "Issue",
	}
	projectCardJSON, _ := json.Marshal(projectCard)
	_, err := g.githubRequest("POST", fmt.Sprintf("/repos/%s/%s/projects/columns/%d/cards", g.Owner, g.Repo, column.ID), string(projectCardJSON))
	if err != nil {
		log.Fatal(err)
	}
}
