package ewelink

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	ewelink, mux, teardown := setupWithHTTP()
	defer teardown()

	mux.HandleFunc("/", fileToResponseHandler(t, "resource/test_data/authentication_ok_response.json"))

	session, err := ewelink.Authenticate(
		context.Background(), NewConfiguration("us"), NewEmailAuthenticator("user@gmail.com", "secret password"))
	failOnError(t, err)

	if session == nil {
		t.Fatal("Got nil session")
	}
}

func TestAuthenticationWithEmail(t *testing.T) {
	ewelink, mux, teardown := setupWithHTTP()
	defer teardown()

	mux.HandleFunc("/", fileToResponseHandler(t, "resource/test_data/authentication_ok_response.json"))

	session, err := ewelink.AuthenticateWithEmail(
		context.Background(), NewConfiguration("us"), "user@gmail.com", "secret password")
	failOnError(t, err)

	if session == nil {
		t.Fatal("Got nil session")
	}
}

func TestAuthenticationWithPhoneNumber(t *testing.T) {
	ewelink := New()

	// Implement once the phone number authenticator has been implemented
	assert.Panics(t, func() {
		// nolint:errcheck
		_, _ = ewelink.AuthenticateWithPhoneNumber(context.Background(), NewConfiguration("us"), "phone number", "secret password")
	})
}

func TestGetDevices(t *testing.T) {
	ewelink, mux, teardown := setupWithHTTP()
	defer teardown()

	mux.HandleFunc("/", fileToResponseHandler(t, "resource/test_data/authentication_ok_response.json"))

	response, err := ewelink.GetDevices(context.Background(), mockSession())
	failOnError(t, err)

	if response.Devicelist == nil {
		t.Fatal("Got nil device list")
	}

	// todo, validate some field in devices response
}

func TestGetDevice(t *testing.T) {
	ewelink, _, teardown := setupWithHTTP()
	defer teardown()

	response, err := ewelink.GetDevice(context.Background(), mockSession(), "deviceId")
	failOnError(t, err)

	if response == nil {
		t.Fatal("Got nil response")
	}

	// TODO, validate some field in device response
}

func TestSetDevicePowerState(t *testing.T) {
	ewelink, mux, address, teardown := setupWithWebsocket()
	defer teardown()

	session := mockSession()
	session.Configuration.WebsocketURL = &url.URL{Scheme: "ws", Host: address.String(), Path: websocketURI}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}

		_ = c.WriteMessage(1, []byte("{\"error\":0,\"apikey\":\"394546b8-6ee3-4f7d-b997-39f057e868b0\",\"config\":{\"hb\":1,\"hbInterval\":145},\"sequence\":\"1\"}"))
		_ = c.WriteMessage(1, []byte("{\"error\":0,\"deviceid\":\"10002c27f2\",\"apikey\":\"394546b8-6ee3-4f7d-b997-39f057e868b0\",\"sequence\":\"2\"}\n"))
	})

	response, err := ewelink.SetDevicePowerState(context.Background(), session, &Device{Uiid: 1, DeviceID: ""}, false)
	failOnError(t, err)

	if response == nil {
		t.Fatal("Got nil response")
	}
}




// setupWithHttp sets up a test http server along with ewelink.
func setupWithHTTP() (*Ewelink, *http.ServeMux, func()) {
	httpClient, mux, teardown := setupHTTPClientAndServer()

	ewelink := New(withHTTPClient(httpClient))

	return ewelink, mux, teardown
}

func setupWithWebsocket() (*Ewelink, *http.ServeMux, net.Addr, func()) {
	ewelink := New()
	mux, teardown, address := setupWebsocket()

	return ewelink, mux, address, teardown
}

func mockSession() *Session {
	return &Session{
		Application:         newApplication(),
		User:                &User{AppID: "1", APIKey: "1", Language: "en"},
		AuthenticationToken: "",
		Configuration:       NewConfiguration("us"),
		MobileDevice:        newIOSDevice(),
	}
}

func fileToResponseHandler(t *testing.T, filename string) http.HandlerFunc {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		// nolint:errcheck
		// #nosec G104
		io.Copy(w, file)
		// nolint:errcheck
		// #nosec G104
		file.Close()
	}
}

func setupHTTPClientAndServer() (*http.Client, *http.ServeMux, func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", http.StripPrefix("/", mux))

	s := httptest.NewTLSServer(mux)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			// #nosec G402
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, mux, s.Close
}

func setupWebsocket() (*http.ServeMux, func(), net.Addr) {
	mux := http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", http.StripPrefix("/", mux))

	s := httptest.NewServer(mux)

	return mux, s.Close, s.Listener.Addr()
}

func failOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
