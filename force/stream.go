package force

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type forceStreaming struct {
	ClientId             string
	SubscribedPushTopics map[string]func(...interface{})
	Timeout              int
	forceApi             *ForceApi
	longPollClient       *http.Client
}

func (s *forceStreaming) httpPost(payload string) (*http.Response, error) {
	ioPayload := strings.NewReader(payload)
	endpoint := s.forceApi.oauth.InstanceUrl + "/cometd/33.0" //version needs to be dynamic
	headerVal := "OAuth " + s.forceApi.oauth.AccessToken

	request, _ := http.NewRequest("POST", endpoint, ioPayload)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", headerVal)

	resp, err := s.longPollClient.Do(request)

	return resp, err
}

func (s *forceStreaming) connect() ([]byte, error) {
	connectParams := `{ "channel": "/meta/connect", "clientId": "` + s.ClientId + `", "connectionType": "long-polling"}`
	resp, err := s.httpPost(connectParams)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return respBytes, err
}

func (forceApi *ForceApi) ConnectToStreamingApi() {
	//set up the client
	cookiejarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookiejarOptions)
	forceApi.stream = &forceStreaming{"", nil, 0, forceApi, &http.Client{Jar: jar}}

	//handshake
	var params = `{"channel":"/meta/handshake", "supportedConnectionTypes":["long-polling"], "version":"1.0"}`
	handshakeResp, _ := forceApi.stream.httpPost(params)
	handshakeBytes, _ := ioutil.ReadAll(handshakeResp.Body)
	defer handshakeResp.Body.Close()

	var data []map[string]interface{}
	json.Unmarshal(handshakeBytes, &data)
	fmt.Println(data)
	forceApi.stream.ClientId = data[0]["clientId"].(string)

	//must handle error here

	// connect
	connResp, _ := forceApi.stream.connect()
	fmt.Println(string(connResp))
	go func() {
		// got to allow disconnect, handle errors
		for {
			connResp, _ = forceApi.stream.connect()
			fmt.Println(string(connResp))
		}
	}()
}

//here we have to allow the ability to pass in a callback function
func (forceApi *ForceApi) SubscribeToPushTopic(pushTopic string) {
	topicString := "/topic/" + pushTopic
	subscribeParams := `{ "channel": "/meta/subscribe", "clientId": "` + forceApi.stream.ClientId + `", "subscription": "` + topicString + `"}`
	subscribeResp, _ := forceApi.stream.httpPost(subscribeParams)
	subscribeBytes, _ := ioutil.ReadAll(subscribeResp.Body)
	defer subscribeResp.Body.Close()
	fmt.Println(string(subscribeBytes))
}

func UnsubscribeFromPushTopic(pushTopic string) {
	fmt.Println(pushTopic)
}

func DisconnectStreamingApi() {
}
