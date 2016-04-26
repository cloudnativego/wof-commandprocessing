package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudnativego/wof-eventprocessing/events"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	request   *http.Request
	recorder  *httptest.ResponseRecorder
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func makeTestServer(dispatchers dispatcherMap) *negroni.Negroni {
	server := negroni.New() // don't need all the middleware here or logging.
	mx := mux.NewRouter()
	initRoutes(mx, formatter, dispatchers)
	server.UseHandler(mx)
	return server
}

func TestAddMoveDispatchesMoveEvent(t *testing.T) {
	dispatchers := make(dispatcherMap)
	dispatcher := newFakeQueueDispatcher()
	dispatchers[MovesQueueName] = dispatcher

	server := makeTestServer(dispatchers)
	recorder := httptest.NewRecorder()
	body := []byte("{\"player_id\":\"D'Pez\", \"target_tile_id\": \"tile-guid\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/game-id/moves", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new move event to return 201, got %d", recorder.Code)
	}
	if len(dispatcher.Messages) != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", len(dispatcher.Messages))
	}
	event := dispatcher.Messages[0].(events.PlayerMovedEvent)
	if event.GameID != "game-id" {
		t.Errorf("Error retrieving GameID from URL. Expected 'game-id', received %s", event.GameID)
	}
	if event.Timestamp <= 0 || event.Timestamp > time.Now().Unix() {
		t.Errorf("Event timestamp was not correctly set. Received: %d", event.Timestamp)
	}
	if event.PlayerID != "D'Pez" {
		t.Errorf("Expected PlayerID to equal \"D'Pez\"; received %s", event.PlayerID)
	}
	if event.TargetTileID != "tile-guid" {
		t.Errorf("Expected TargetTileID to equal 'tile-guid'; received %s", event.TargetTileID)
	}
}

func TestPlayerJoinDispatchesPlayerJoinEvent(t *testing.T) {
	dispatchers := make(dispatcherMap)
	dispatcher := newFakeQueueDispatcher()
	dispatchers[PlayerJoinsQueueName] = dispatcher

	server := makeTestServer(dispatchers)
	recorder := httptest.NewRecorder()
	body := []byte("{\"player_id\":\"Swordless\", \"sprite\": \"mime\", \"name\": \"Mimetown\"}")
	reader := bytes.NewReader(body)
	request, _ = http.NewRequest("POST", "/game-id/join", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected creation of new player join event to return 201, got %d", recorder.Code)
	}
	if len(dispatcher.Messages) != 1 {
		t.Errorf("Expected queue dispatch count of 1, got %d", len(dispatcher.Messages))
	}
	event := dispatcher.Messages[0].(events.PlayerJoinedEvent)
	if event.GameID != "game-id" {
		t.Errorf("Error retrieving GameID from URL. Expected 'game-id', received %s", event.GameID)
	}
	if event.Timestamp <= 0 || event.Timestamp > time.Now().Unix() {
		t.Errorf("Event timestamp was not correctly set. Received: %d", event.Timestamp)
	}
	if event.PlayerID != "Swordless" {
		t.Errorf("Expected PlayerID to equal 'Swordless'; received %s", event.PlayerID)
	}
	if event.Sprite != "mime" {
		t.Errorf("Expected Sprite to equal 'mime'; received %s", event.Sprite)
	}
	if event.Name != "Mimetown" {
		t.Errorf("Expected Name to equal 'Mimetown'; received %s", event.Name)
	}
}
