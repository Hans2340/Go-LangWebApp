package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Post struct to hold the post data
type UserRetrieve struct {
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

type Postings struct {
	User_ID string `json:"user_id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
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
	var UserRetrieve []UserRetrieve
	err = json.Unmarshal(body, &UserRetrieve)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Handle / route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, UserRetrieve)

	})
	http.HandleFunc("/create:user", createUser)
	http.HandleFunc("/create:post", createPostings)
	http.HandleFunc("/delete:user", deleteUser)
	http.HandleFunc("/update", updateUser)

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
func createPostings(w http.ResponseWriter, r *http.Request) {

	// parse the post data from the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// create a new user struct and populate it with the form data
	posting := &Postings{
		User_ID: r.FormValue("user_id"),
		Title:   r.FormValue("title"),
		Body:    r.FormValue("body"),
	}

	// convert the user struct to json
	postingData, err := json.Marshal(posting)
	if err != nil {
		http.Error(w, "Error marshalling user data", http.StatusInternalServerError)
		return
	}

	// create a new request to the API
	req, err := http.NewRequest("POST", "https://gorest.co.in/public/v2/posts", bytes.NewBuffer(postingData))
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
func updateUser(w http.ResponseWriter, r *http.Request) {

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

	// get the id of the user to be updated from the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "User id is required", http.StatusBadRequest)
		return
	}

	// create a new request to the API
	req, err := http.NewRequest("PUT", "https://gorest.co.in/public/v2/users/"+id, bytes.NewBuffer(userData))
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
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the user id from the URL
	userId := mux.Vars(r)["userId"]

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new DELETE request with the user id
	req, err := http.NewRequest("DELETE", "https://gorest.co.in/public/v2/users/"+userId, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the DELETE request
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response body to the HTTP response
	w.Write(body)
}
