package goss_go_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (connect *ConnectionConfig) GenerateEndpoint(target string) string {
	var sb strings.Builder

	if connect.https {
		sb.WriteString("https://")
	} else {
		sb.WriteString("http://")
	}
	sb.WriteString(connect.host + ":")
	sb.WriteString(strconv.Itoa(connect.port))
	sb.WriteString(target)
	return sb.String()
}

const (
	createTarget   = "/create"
	getTarget      = "/get"
	deleteTarget   = "/delete"
	truncateTarget = "/truncate"
)

type ConnectionConfig struct {
	host  string
	port  int
	https bool
}

type SessionData struct {
	Session string `json:"session"`
	TTL     int    `json:"ttl"`
}

type GeneralResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (connect *ConnectionConfig) CreateSession(session string, ttl int) (string, error) {
	// converging session with ttl for json
	sessionData := SessionData{
		Session: session,
		TTL:     ttl,
	}
	// struct to json
	sessionJson, err := json.Marshal(sessionData)
	fmt.Print(string(sessionJson))
	if err != nil {
		return "", errors.New("invalid data for session")
	}

	// getting the url
	url := connect.GenerateEndpoint(createTarget)

	// forming the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(sessionJson)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	// making the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// closing the body
	defer resp.Body.Close()

	// reading the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// extracting general response with status and message
	var generalResponse GeneralResponse
	err = json.Unmarshal(body, &generalResponse)
	if err != nil {
		return "", err
	}

	// checking if status 201 for created
	if generalResponse.Status != 201 {
		return "", errors.New("session not created")
	}

	// returning session id
	return generalResponse.Message, nil
}

func (connect *ConnectionConfig) GetSession(sessionId string) (string, error) {

	if len(sessionId) == 0 {
		return "", errors.New("invalid session id")
	}

	// getting the url
	url := connect.GenerateEndpoint(getTarget) + "?sessionId=" + sessionId
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// closing the body
	defer res.Body.Close()

	// reading the body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// checking if the response is 200
	if res.Status != "200 OK" {
		return "", errors.New("status not found")
	}

	// returning the session
	return string(body), nil
}
func (connect *ConnectionConfig) DeleteSession(sessionId string) error {

	if len(sessionId) == 0 {
		return errors.New("invalid session id")
	}

	// getting the url
	url := connect.GenerateEndpoint(deleteTarget) + "?sessionId=" + sessionId
	// forming the request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// making the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	// checking if the response is 200
	if res.Status != "200 OK" {
		return errors.New("session not found")
	}
	// returning the session
	return nil
}
func (connect *ConnectionConfig) TruncateStore() error {
	// getting the url
	url := connect.GenerateEndpoint(truncateTarget)
	// forming the request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// making the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	// checking if the response is 200
	if res.Status != "200 OK" {
		return errors.New("could not truncate store")
	}
	// returning the session
	return nil
}
