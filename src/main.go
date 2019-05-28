package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

/*
Lock represents a simple lock structure
*/
type Lock struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Workflow   string    `json:"workflow"`
	Created    time.Time `json:"created"`
	LastChange time.Time `json:"lastChange"`
}

var endpointPtr = flag.String("api-endpoint", os.Getenv("API_ENDPOINT"), "The API endpoint URL (defaults to $API_ENDPOINT env variable)")

func createLock(workflow string) (string, error) {

	url, err := url.Parse(*endpointPtr)
	if err != nil {
		log.Fatal("URLMalformated: ", err)
		return "", err
	}

	locJSON, err := json.Marshal(Lock{Workflow: workflow})

	// Build the request
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(locJSON))
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var lock Lock
		if err := json.NewDecoder(resp.Body).Decode(&lock); err != nil {
			log.Fatal("ParseGetResponse: ", err)
			return "", err
		}

		return lock.ID, nil
	}

	return "", fmt.Errorf("API returned status code %d", resp.StatusCode)

}

func crudLock(id string, method string) (string, error) {

	url, err := url.Parse(*endpointPtr)
	if err != nil {
		log.Fatal("URLMalformated: ", err)
		return "", err
	}

	q := url.Query()
	q.Set("id", id)
	url.RawQuery = q.Encode()

	// Build the request
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		return string(body), err
	}

	return "", fmt.Errorf("API returned status code %d", resp.StatusCode)

}

func releaseLock(id string) error {
	_, err := crudLock(id, "PATCH")
	return err
}

func deleteLock(id string) error {
	_, err := crudLock(id, "DELETE")
	return err
}

func getLock(id string) (Lock, error) {

	var lock Lock

	body, err := crudLock(id, "GET")
	if err != nil {
		return lock, err
	}

	err = json.Unmarshal([]byte(body), &lock)
	return lock, err
}

func printUsage() {
	fmt.Println("lock-manager [--api-endpoint api-url] {create workflow|release id|delete id}")
}

func main() {

	// Parse flags
	flag.Parse()

	if *endpointPtr == "" {
		log.Fatal("API endpoint not defined (set environment variable API_ENDPOINT or provide --api-endpoint flag)")
	}

	// Parse command
	if len(flag.Args()) < 2 {
		fmt.Println("ERROR: Insufficient number of arguments.")
		printUsage()
		os.Exit(1)
	}

	cmd := flag.Arg(0)

	var err error
	var id string

	switch cmd {

	case "create":
		workflow := flag.Arg(1)
		id, err = createLock(workflow)
		fmt.Println(id)
		break

	case "get":
		id := flag.Arg(1)
		lock, err := getLock(id)
		if err == nil {
			locJSON, _ := json.Marshal(lock)
			fmt.Print(string(locJSON))
		}
		break

	case "release":
		id := flag.Arg(1)
		err = releaseLock(id)
		break

	case "delete":
		id := flag.Arg(1)
		err = deleteLock(id)
		break

	default:
		fmt.Printf("Unknown command [%s]", flag.Arg(0))
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
