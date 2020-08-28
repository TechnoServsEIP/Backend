package controllers

import "net/http"

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := GetGithubID()
	redirectURL :=
}
