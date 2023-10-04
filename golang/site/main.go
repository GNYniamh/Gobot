package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageVariables struct {
	Title string
}

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/login-submit", LoginSubmit)

	fmt.Println("Server is listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Home Page",
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Login Page",
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// You can add your authentication logic here.
		// For this example, we will just print the received values.
		fmt.Printf("Received login request. Username: %s, Password: %s\n", username, password)

		// Redirect to the home page after login.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}
