package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// Post struct to hold the post data
type Post struct {
	ID     float64 `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Gender string  `json:"gender"`
	Status string  `json:"status"`
}

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
	Status string `json:"status"`
}

func main() {
	// Parse the HTML template
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Make GET request to API
	res, err := http.Get("https://gorest.co.in/public/v2/users")
	if err != nil {
		fmt.Println("Error fetching data from API:", err)
		return
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal JSON response into Post struct
	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Handle / route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, posts)
	})
	http.HandleFunc("/create", createUser)

	// Start server
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// //////////////////////////////////////////////////////////////////////
// creating new users
func createUser(w http.ResponseWriter, r *http.Request) {

	// parse the post data from the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// create a new user struct and populate it with the form data
	user := &User{
		Name:   r.FormValue("name"),
		Email:  r.FormValue("email"),
		Gender: r.FormValue("gender"),
		Status: r.FormValue("status"),
	}

	// convert the user struct to json
	userData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshalling user data", http.StatusInternalServerError)
		return
	}

	// create a new request to the API
	req, err := http.NewRequest("POST", "https://gorest.co.in/public/v2/users", bytes.NewBuffer(userData))
	if err != nil {
		http.Error(w, "Error creating new request", http.StatusInternalServerError)
		return
	}

	// add the necessary headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer acb7f761573e6f253cc1fd40f95cad4a20bf9a1ca4b1f2f71309ec190746d5fb")

	// make the request to the API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request to API", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}
		bodyString := string(bodyBytes)
		fmt.Println("API response body:", bodyString)
	}
}
