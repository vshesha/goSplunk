package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {

	//User inputs
	username := flag.String("u", "", "Splunk username")
	password := flag.String("p", "", "Splunk password")
	path := flag.String("f", "", "Splunk app file path")

	flag.Parse()
	if len(*username) == 0 || len(*password) == 0 || len(*path) == 0 {
		fmt.Println("Missing params. See help (-h).")
		os.Exit(0)
	}

	//Login User
	usertoken := loginUser(*username, *password)

	//Submit App
	requestId := submitApp(*path, usertoken)

	//Check status
	temp := time.Now()
	for {
		status := statusCheck(requestId, usertoken)
		if status == "ERROR" {
			break
		}
		if status == "SUCCESS" {
			break
		}
		fmt.Println("Sleeping 30 seconds ... ")
		time.Sleep(30 * time.Second)
		if time.Since(temp) > 5*time.Minute {
			//Watchdog
			fmt.Println("Waited too long")
			os.Exit(1)
		}
	}

	//Report Status
	results := statusReport(requestId, usertoken)
	fmt.Println(results)
}

func loginUser(username, password string) string {

	req, err := http.NewRequest("GET", "https://api.splunk.com/2.0/rest/login/splunk", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Generating User Auth Token")
	return responseToJSON(resp).Get("data").Get("token").MustString()
}

func submitApp(path, usertoken string) string {

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	// Create file field
	fw, err := w.CreateFormFile("app_package", path)
	if err != nil {
		log.Fatalln(err)
	}
	fd, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatalln(err)
	}
	// terminating boundary
	w.Close()

	req, err := http.NewRequest("POST", "https://appinspect.splunk.com/v1/app/validate", buf)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "Bearer "+usertoken)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return responseToJSON(resp).Get("request_id").MustString()
}

func statusCheck(requestId, usertoken string) string {
	req, err := http.NewRequest("GET", "https://appinspect.splunk.com/v1/app/validate/status/"+requestId, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+usertoken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return responseToJSON(resp).Get("status").MustString()
}

func statusReport(requestId, usertoken string) string {
	req, err := http.NewRequest("GET", "https://appinspect.splunk.com/v1/app/report/"+requestId, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+usertoken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Println(responseString)
	return responseString
}

func responseToJSON(resp *http.Response) *simplejson.Json {
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	fmt.Println(responseString)

	js, err := simplejson.NewJson(responseData)
	if err != nil {
		log.Fatalln(err)
	}
	return js
}
