package api

import (
	"net/http"
)

func setupRoutes() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/deadlines", deadlinesHandler)
	http.HandleFunc("/settings", settingsHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the homepage.")
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the tasks endpoint.")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the users endpoint.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the login endpoint.")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the register endpoint.")
}

func deadlinesHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the deadlines endpoint.")
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	print(w, "This is the settings endpoint.")
}
