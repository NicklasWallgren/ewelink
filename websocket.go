package ewelink

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	websocketScheme = "wss"
	websocketURI    = "/api/ws"
)

// WebsocketClient is the interface implemented by types that can invoke the Ewelink Websocket Api.
type WebsocketClient interface {
	call(context context.Context, requests []WebsocketRequest, session *Session) ([]*requestResult, error)
}

type websocketClient struct {
	encoder encoder
	decoder responseDecoder
}

func newWebsocketClient() *websocketClient {
	return &websocketClient{encoder: newJSONEncoder(), decoder: newJSONResponseDecoder()}
}

type requestResult struct {
	Response Response
	Error    error
}

func (w websocketClient) connect(context context.Context, url *url.URL) (*websocket.Conn, error) {
	c, response, err := websocket.DefaultDialer.DialContext(context, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to establish websocket connection %w", err)
	}

	defer response.Body.Close() // nolint:errcheck

	return c, nil
}

func (w websocketClient) call(context context.Context, requests []WebsocketRequest, session *Session) ([]*requestResult, error) {
	connection, err := w.connect(context, session.Configuration.WebsocketURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect. %w", err)
	}

	defer connection.Close() // nolint:errcheck

	return w.do(connection, requests), nil
}

func (w websocketClient) do(connection *websocket.Conn, requests []WebsocketRequest) []*requestResult {
	responses := make([]*requestResult, len(requests))

	for i, request := range requests {
		response, err := w.request(connection, request)

		responses[i] = &requestResult{Response: response, Error: err}
	}

	return responses
}

func (w websocketClient) request(connection *websocket.Conn, request WebsocketRequest) (Response, error) {
	if err := w.sendMessage(connection, request.Payload()); err != nil {
		return nil, err
	}

	return w.readMessage(connection, request.Response())
}

func (w websocketClient) sendMessage(connection *websocket.Conn, payload payload) error {
	encoded, err := w.encoder.encode(payload)
	if err != nil {
		return fmt.Errorf("unable to send websocket message %w", err)
	}

	return connection.WriteMessage(websocket.TextMessage, encoded)
}

func (w websocketClient) readMessage(connection *websocket.Conn, response Response) (Response, error) {
	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("unable to read message from websocket connection %w", err)
	}

	return w.decoder.decode(response, ioutil.NopCloser(bytes.NewReader(message)), http.StatusOK)
}
