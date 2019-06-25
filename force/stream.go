package force

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type forceStreaming struct {
	ClientID       string
	Subscriptions  map[string]func([]byte, ...interface{})
	Timeout        int
	forceAPI       *API
	longPollClient *http.Client
}

func (s *forceStreaming) httpPost(payload string) (*http.Response, error) {
	ioPayload := strings.NewReader(payload)
	endpoint := s.forceAPI.oauth.InstanceURL + "/cometd/33.0" //version needs to be dynamic
	headerVal := "OAuth " + s.forceAPI.oauth.AccessToken

	request, _ := http.NewRequest("POST", endpoint, ioPayload)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", headerVal)

	resp, err := s.longPollClient.Do(request)

	return resp, err
}

func (s *forceStreaming) connect() ([]byte, error) {
	connectParams := `{ "channel": "/meta/connect", "clientId": "` + s.ClientID + `", "connectionType": "long-polling"}`
	resp, err := s.httpPost(connectParams)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return respBytes, err
}

// ConnectToStreamingAPI connects to streaming API prior to issuing subscriptions
func (forceAPI *API) ConnectToStreamingAPI() {
	//set up the client
	cookiejarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookiejarOptions)
	forceAPI.stream = &forceStreaming{"", map[string]func([]byte, ...interface{}){}, 0, forceAPI, &http.Client{Jar: jar}}

	//handshake
	var params = `{"channel":"/meta/handshake", "supportedConnectionTypes":["long-polling"], "version":"1.0"}`
	handshakeResp, _ := forceAPI.stream.httpPost(params)
	handshakeBytes, _ := ioutil.ReadAll(handshakeResp.Body)
	defer handshakeResp.Body.Close()

	var data []map[string]interface{}
	json.Unmarshal(handshakeBytes, &data)
	fmt.Println(data)
	forceAPI.stream.ClientID = data[0]["clientId"].(string)

	//must handle error here

	// connect
	connBytes, _ := forceAPI.stream.connect()

	var connectData []map[string]interface{}
	json.Unmarshal(connBytes, &connectData)
	for _, msg := range data {
		cb := forceAPI.stream.Subscriptions[msg["channel"].(string)]
		if cb != nil {
			cb(connBytes)
		}
		fmt.Println(string(connBytes))
	}

	go func() {
		// got to allow disconnect, handle errors
		for {
			connBytes, _ = forceAPI.stream.connect()
			json.Unmarshal(connBytes, &connectData)

			for _, msg := range connectData {
				cb := forceAPI.stream.Subscriptions[msg["channel"].(string)]
				if cb != nil {
					cb(connBytes)
				}
			}
			//fmt.Println(string(connBytes))
		}
	}()
}

// SubscribeToPushTopic here we have to allow the ability to pass in a callback function
func (forceAPI *API) SubscribeToPushTopic(pushTopic string, callback func([]byte, ...interface{})) ([]byte, error) {
	topicString := "/topic/" + pushTopic
	subscribeParams := `{ "channel": "/meta/subscribe", "clientID": "` + forceAPI.stream.ClientID + `", "subscription": "` + topicString + `"}`

	subscribeResp, _ := forceAPI.stream.httpPost(subscribeParams)
	subscribeBytes, err := ioutil.ReadAll(subscribeResp.Body)

	defer subscribeResp.Body.Close()

	forceAPI.stream.Subscriptions[topicString] = callback
	return subscribeBytes, err

}

// SubscribeToEvent here we have to allow the ability to pass in a callback function
func (forceAPI *API) SubscribeToEvent(eventName string, callback func([]byte, ...interface{})) ([]byte, error) {
	eventString := "/event/" + eventName
	subscribeParams := `{ "channel": "/meta/subscribe", "clientID": "` + forceAPI.stream.ClientID + `", "subscription": "` + eventString + `"}`

	subscribeResp, _ := forceAPI.stream.httpPost(subscribeParams)
	subscribeBytes, err := ioutil.ReadAll(subscribeResp.Body)

	defer subscribeResp.Body.Close()

	forceAPI.stream.Subscriptions[eventString] = callback
	return subscribeBytes, err

}

// UnsubscribeFromPushTopic is a func that doesn't do anything yet
func UnsubscribeFromPushTopic(pushTopic string) {
	fmt.Println(pushTopic)
}

// DisconnectStreamingAPI is a func that doesn't do anything yet
func DisconnectStreamingAPI() {
}
