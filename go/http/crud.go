package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type User struct {
	Role       string `json:"role"`
	ID         string `json:"id"`
	Experience int    `json:"experience"`
	Remote     bool   `json:"remote"`
	User       struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Age      int    `json:"age"`
	} `json:"user"`
}

func getUsers(url string) ([]User, error) {
	// get response
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	defer resp.Body.Close()

	// Unmarshal or decode if bigger data set
	var data []User
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return data, nil
}

func createUser(url, apiKey string, data User) (User, error) {
	// encode data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return User{}, err
	}

	// create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return User{}, err
	}

	// set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// create a new client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	// decode the json data from the response into a new User struct
	var u User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func updateUser(baseURL, id, apiKey string, data User) (User, error) {
	fullURL := baseURL + "/" + id

	// encode data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return User{}, err
	}

	// create new request
	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return User{}, err
	}

	// set content type and api key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// create client and Do Put
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}

	// decode client response and return user
	var u User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func getUserById(baseURL, id, apiKey string) (User, error) {
	fullURL := baseURL + "/" + id

	fmt.Println("______")
	b := []byte{}
	fmt.Println(b)
	fmt.Println("______")

	// create new request, set body as nil, not empty []byte{}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return User{}, err
	}
	// set content type and api key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// create client and Do Put
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}

	// decode client response and return user
	var u User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func deleteUser(baseURL, id, apiKey string) error {
	fullURL := baseURL + "/" + id
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return err
	}
	// set content type and api key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	// create client and Do Put
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// check the statusCode
	if resp.StatusCode > 200 {
		return err
	}

	return nil
}

type Issue struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Estimate int    `json:"estimate"`
	Status   string `json:"status"`
}

func getIssues(url string) []Issue {
	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var issues []Issue
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&issues)
	if err != nil {
		return nil
	}

	return issues
}

func logIssues(issues []Issue) string {
	log := ""
	for _, issue := range issues {
		log += fmt.Sprintf("- Issue: %s - Estimate: %d\n", issue.Title, issue.Estimate)
	}
	return log
}

func fetchTasks(baseURL, availability string) []Issue {

	fullURL := baseURL + "?sort=estimate"

	var limit int
	switch availability {
	case "Low":
		limit = 1
	case "Medium":
		limit = 3
	case "High":
		limit = 5
	}

	fullURL += "&limit=" + strconv.Itoa(limit)

	return getIssues(fullURL)

}
