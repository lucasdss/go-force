package force

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	connectParams := `{ "channel": "/meta/connect", "clientID": "` + s.ClientID + `", "connectionType": "long-polling"}`
	resp, err := s.httpPost(connectParams)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return respBytes, err
}

// UnsubscribeFromPushTopic is a func that doesn't do anything yet
func UnsubscribeFromPushTopic(pushTopic string) {
	fmt.Println(pushTopic)
}

// DisconnectStreamingAPI is a func that doesn't do anything yet
func DisconnectStreamingAPI() {
}
