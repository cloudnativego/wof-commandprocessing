package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func makeTestServer() *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	repo := NewFakeRepository()
	initRoutes(mx, formatter, repo)
	server.UseHandler(mx)
	return server
}

func TestGetForNonexistantMapReturns404(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := makeTestServer()
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/api/maps/nevergonnahappen", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected a 404 for non-existent map, got %d instead.", recorder.Code)
	}
}

func TestGetMapList(t *testing.T) {
	server := makeTestServer()

	for x := 0; x < 10; x++ {
		newMap := generateTestMap(10, "Testosthenes")
		mapBytes, _ := json.Marshal(newMap)
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("PUT", fmt.Sprintf("/api/maps/%s", newMap.ID), bytes.NewBuffer(mapBytes))
		server.ServeHTTP(recorder, request)
	}

	recorder2 := httptest.NewRecorder()
	request2, _ := http.NewRequest("GET", "/api/maps", nil)
	server.ServeHTTP(recorder2, request2)
	var mapResponse []WofMap
	err := json.Unmarshal(recorder2.Body.Bytes(), &mapResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response from server: %s", err.Error())
		return
	}

	if len(mapResponse) != 10 {
		t.Errorf("Expected 10 maps, got %d", len(mapResponse))
		return
	}

	if mapResponse[0].Metadata.Author != "Testosthenes" {
		t.Errorf("Expected Testosthenes as author, got '%s'", mapResponse[0].Metadata.Author)
	}

}

func TestUpdateAndRetrieveMap(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	server := makeTestServer()
	recorder = httptest.NewRecorder()
	testMap := generateTestMap(10, "Testophocles")
	mapBytes, err := json.Marshal(testMap)
	request, _ = http.NewRequest("PUT", "/api/maps/testmap", bytes.NewBuffer(mapBytes))
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected %v; received %v", http.StatusCreated, recorder.Code)
	}

	recorder2 := httptest.NewRecorder()
	request2, _ := http.NewRequest("GET", "/api/maps/testmap", nil)
	server.ServeHTTP(recorder2, request2)

	var mapResponse WofMap
	err = json.Unmarshal(recorder2.Body.Bytes(), &mapResponse)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %s", err)
		return
	}

	if mapResponse.Metadata.Author != "Testophocles" {
		t.Errorf("Expected `Testophocles` author on updated map, got '%s'", mapResponse.Metadata.Author)
		return
	}

	testMap.Tiles[0][0].TileName = "updated-tile"
	testMap.Tiles[0][0].Sprite = "torch"
	mapBytes, err = json.Marshal(testMap)

	recorder3 := httptest.NewRecorder()
	request3, _ := http.NewRequest("PUT", "/api/maps/testmap", bytes.NewBuffer(mapBytes))
	server.ServeHTTP(recorder3, request3)

	recorder4 := httptest.NewRecorder()
	request4, _ := http.NewRequest("GET", "/api/maps/testmap", nil)
	server.ServeHTTP(recorder4, request4)
	var mapResponse2 WofMap
	err = json.Unmarshal(recorder4.Body.Bytes(), &mapResponse2)
	if err != nil {
		t.Errorf("Error unmarshaling map details (2nd query): %s", err)
		return
	}

	if mapResponse2.Tiles[0][0].TileName != "updated-tile" {
		t.Errorf("Expected 'updated-tile' as tile name, got '%s'", mapResponse2.Tiles[0][0].TileName)
	}

	if mapResponse2.ID != mapResponse.ID {
		t.Errorf("Expected IDs to remain the same, were different. Got %s and %s", mapResponse2.ID, mapResponse.ID)
	}
}

func generateTestMap(size int, author string) (gameMap WofMap) {

	gameMap.Metadata.Author = author
	gameMap.Metadata.Description = "Auto-generated Test Map"
	gameMap.ID = uuid.New()

	tiles := make([][]MapTile, size)
	for row := 0; row < size; row++ {
		tiles[row] = make([]MapTile, size)
		for col := 0; col < size; col++ {
			tiles[row][col] = makeTile()
		}
	}
	gameMap.Tiles = tiles

	return
}

func makeTile() (tile MapTile) {
	tile.AllowDown = true
	tile.AllowLeft = true
	tile.AllowRight = true
	tile.AllowUp = true
	tile.Sprite = ""
	tile.TileName = "test-tile"
	tile.ID = uuid.New()
	return
}
