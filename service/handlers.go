package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloudnativego/wof-eventprocessing/events"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func addMoveHandler(formatter *render.Render, dispatcher QueueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newMoveCommand moveCommand
		err := json.Unmarshal(payload, &newMoveCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse add move command.")
			return
		}

		vars := mux.Vars(req)
		gameID := vars["gameID"]

		evt := events.PlayerMovedEvent{
			GameID:       gameID,
			PlayerID:     newMoveCommand.PlayerID,
			TargetTileID: newMoveCommand.TargetTileID,
			Timestamp:    time.Now().Unix(),
		}

		dispatcher.DispatchMessage(evt)
		formatter.JSON(w, http.StatusCreated, evt)
	}
}

func playerJoinHandler(formatter *render.Render, dispatcher QueueDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newPlayerJoinCommand playerJoinCommand
		err := json.Unmarshal(payload, &newPlayerJoinCommand)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse player join command.")
			return
		}

		vars := mux.Vars(req)
		gameID := vars["gameID"]

		evt := events.PlayerJoinedEvent{
			GameID:    gameID,
			PlayerID:  newPlayerJoinCommand.PlayerID,
			Sprite:    newPlayerJoinCommand.Sprite,
			Name:      newPlayerJoinCommand.Name,
			Timestamp: time.Now().Unix(),
		}

		dispatcher.DispatchMessage(evt)
		formatter.JSON(w, http.StatusCreated, evt)
	}
}
