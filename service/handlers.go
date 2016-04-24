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
