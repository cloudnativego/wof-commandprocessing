package service

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer(appEnv *cfenv.App, dispatchers dispatcherMap) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, dispatchers)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, dispatchers dispatcherMap) {

	mx.HandleFunc("/{gameID}/moves", addMoveHandler(formatter, dispatchers[MovesQueueName])).Methods("POST")
	mx.HandleFunc("/{gameID}/join", playerJoinHandler(formatter, dispatchers[PlayerJoinsQueueName])).Methods("POST")
}
