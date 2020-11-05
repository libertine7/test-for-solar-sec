package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewVacancy struct {
	Forename    string
	Salarylevel int
	Experience  string
	City        string
}

type Vacancy struct {
	Id          int
	Forename    string
	Salarylevel int
	Experience  string
	City        string
}

func PutData() {
	url := "http://127.0.0.1:8000/vacancy"
	body, err := json.Marshal(NewVacancy{Forename: "Golang Dev", Salarylevel: 20000, Experience: "Very good", City: "London"})

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("editor", "123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if Err := scanner.Err(); Err != nil {
		return
	}
}


func GetData() {
	url := "http://127.0.0.1:8000/vacancy"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("viewer", "123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if Err := scanner.Err(); Err != nil {
		return
	}
}


func GetDataByID() {
	url := "http://127.0.0.1:8000/vacancy/2"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("viewer", "123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if Err := scanner.Err(); Err != nil {
		return
	}
}


func DelDataByID() {
	url := "http://127.0.0.1:8000/vacancy/2"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("editor", "123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if Err := scanner.Err(); Err != nil {
		return
	}
}


func main() {
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** Put data")
	PutData()
	fmt.Println("***** get data all")
	GetData()
	fmt.Println("***** GetDataByID")
	GetDataByID()
	fmt.Println("***** DelDataByID")
	DelDataByID()
}
