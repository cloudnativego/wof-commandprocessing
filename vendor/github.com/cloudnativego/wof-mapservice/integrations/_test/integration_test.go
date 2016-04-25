package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.google.com/p/go-uuid/uuid"

	. "github.com/cloudnativego/wof-mapservice/service"
)

var (
	server = NewServer()
)

func TestIntegration(t *testing.T) {

	// Query - empty maps
	emptyMaps, _ := getMapList(t)
	if len(emptyMaps) > 0 {
		t.Errorf("Expected get map list to return an empty array; received %d", len(emptyMaps))
	}

	// Create first map	// Add first match
	newMap := generateTestMap(5, "Mephistopholes")
	mapResponse, _ := putMap(t, newMap)
	if mapResponse.Metadata.Author != "Mephistopholes" {
		t.Errorf("Did not receive identical payload in response to map PUT, got %v", mapResponse)
		return
	}

	// Get Map list again
	updatedMaps, _ := getMapList(t)
	if len(updatedMaps) != 1 {
		t.Errorf("Expected get map list to return 1 item, received %d", len(updatedMaps))
	}

	if updatedMaps[0].Metadata.Author != "Mephistopholes" {
		t.Errorf("First map has the incorrect author. Got '%s'", updatedMaps[0].Metadata.Author)
	}

	// Change the first map
	newMap.Metadata.Author = "Hingle McCringleBerry"
	mapResponse2, _ := putMap(t, newMap)
	if mapResponse2.Metadata.Author != "Hingle McCringleBerry" {
		t.Errorf("Response from PUT did not match request: %v", mapResponse2)
	}

	// Add decoy maps, also test unicode round-tripping
	_, _ = putMap(t, generateTestMap(5, "אֶחָד"))
	_, _ = putMap(t, generateTestMap(5, "שְׁתַיִם"))
	lastMap := generateTestMap(5, "שָׁלוֹשׁ")
	_, _ = putMap(t, lastMap)

	// Get map list (should have 4 now)
	newList, _ := getMapList(t)
	if len(newList) != 4 {
		t.Errorf("Should have 4 maps on file now, instead have %d", len(newList))
	}

	lastMap2, _ := getMap(t, lastMap.ID)
	if lastMap2.Metadata.Author != "שָׁלוֹשׁ" {
		t.Errorf("Unicode round-tripping failed, got %v", lastMap2)
	}

	// Get map details on updated map
	updatedMap, _ := getMap(t, newMap.ID)
	if updatedMap.Metadata.Author != "Hingle McCringleBerry" {
		t.Errorf("Map details were incorrect, got %v", updatedMap)
	}
}

func getMapList(t *testing.T) (maps []WofMap, err error) {
	request, _ := http.NewRequest("GET", "/api/maps", nil)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)
	maps = make([]WofMap, 0)
	err = json.Unmarshal(recorder.Body.Bytes(), &maps)
	if err != nil {
		t.Errorf("Error unmarshaling match list, %v", err)
	} else {
		if recorder.Code != 200 {
			t.Errorf("Expected map list code to be 200, got %d", recorder.Code)
		} else {
			fmt.Printf("Retrieved %d maps.\n", len(maps))
		}
	}

	return
}

func putMap(t *testing.T, gameMap WofMap) (reply WofMap, err error) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(gameMap)
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/api/maps/%s", gameMap.ID), bytes.NewBuffer(body))
	server.ServeHTTP(recorder, request)
	if recorder.Code != 201 {
		t.Errorf("Error creating map, expected 201 code, got %d", recorder.Code)
	}
	var mapResponse WofMap
	err = json.Unmarshal(recorder.Body.Bytes(), &mapResponse)
	if err != nil {
		t.Errorf("Error unmarshaling new map response: %v", err)
	} else {
		fmt.Println("\tAdded Map OK")
	}
	reply = mapResponse
	return
}

func getMap(t *testing.T, ID string) (gameMap WofMap, err error) {
	recorder := httptest.NewRecorder()
	mapURL := fmt.Sprintf("/api/maps/%s", ID)
	request, _ := http.NewRequest("GET", mapURL, nil)
	server.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("Error getting map details: %d", recorder.Code)
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &gameMap)
	if err != nil {
		t.Errorf("Error unmarshaling map details: %v", err)
	} else {
		fmt.Println("\tQueried Map Details OK")
	}
	return
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
