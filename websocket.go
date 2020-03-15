package ewelink

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/url"
)

const websocketScheme = "wss"
const websocketUri = "/api/ws"

type WebsocketClient interface {
	call(context context.Context, requests []WebsocketRequest, session *Session) ([]*result, error)
}

type websocketClient struct {
	encoder       Encoder
	decoder       ResponseDecoder
}

func newWebsocketClient() *websocketClient {
	return &websocketClient{encoder: newJsonEncoder(), decoder: newResponseJsonDecoder()}
}

type result struct {
	Response Response
	Error    error
}

func (w websocketClient) connect(context context.Context, url url.URL) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.DialContext(context, url.String(), nil)

	return c, err
}

func (w websocketClient) call(context context.Context, requests []WebsocketRequest, session *Session) ([]*result, error) {
	connection, err := w.connect(context,
		url.URL{Scheme: websocketScheme, Host: session.Configuration.WebsocketHost, Path: websocketUri})

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	return w.do(connection, requests), nil
}

func (w websocketClient) do(connection *websocket.Conn, requests []WebsocketRequest) []*result {
	var responses = make([]*result, len(requests))

	for i, request := range requests {
		response, err := w.request(connection, request)

		responses[i] = &result{Response: response, Error: err}
	}

	return responses
}

func (w websocketClient) request(connection *websocket.Conn, request WebsocketRequest) (Response, error) {
	if err := w.sendMessage(connection, request.Payload()); err != nil {
		return nil, err
	}

	return w.readMessage(connection, request.Response())
}

func (w websocketClient) sendMessage(connection *websocket.Conn, payload payloadInterface) error {
	encoded, err := w.encoder.encode(payload)

	if err != nil {
		return err
	}

	return connection.WriteMessage(websocket.TextMessage, encoded)
}

func (w websocketClient) readMessage(connection *websocket.Conn, response Response) (Response, error) {
	_, message, err := connection.ReadMessage()

	if err != nil {
		return nil, nil
	}

	fmt.Println(string(message))

	return w.decoder.decode(response, ioutil.NopCloser(bytes.NewReader(message)))
}
