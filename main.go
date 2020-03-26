package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const nplayers = 4

var game *GameState
var playersIn int
var playerTokens []string
var globalLock sync.Mutex

func handleJoin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	globalLock.Lock()
	defer globalLock.Unlock()

	if playersIn >= nplayers {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenData := make([]byte, 16)
	_, err = crand.Read(tokenData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := base64.StdEncoding.EncodeToString(tokenData)

	game.SetName(playersIn, name)
	playersIn = playersIn + 1

	playerTokens = append(playerTokens, token)

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 24 * 3600,
		Path:   "/",
	})
	w.Header().Add("Location", "/play.html")
	w.WriteHeader(http.StatusSeeOther)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, tokenb, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Printf("error: %v\n", err)
		}
		return
	}

	token := string(tokenb)

	playerIdx := -1
	for i, ptoken := range playerTokens {
		if ptoken == token {
			playerIdx = i
		}
	}
	if playerIdx == -1 {
		conn.Close()
		return
	}

	eventStream := make(chan interface{}, 100)
	end := make(chan struct{})
	listener := &GameStateListener{
		l:           eventStream,
		playerIndex: playerIdx,
	}

	game.AddListener(listener)

	// Send events to browser
	go func() {
		for {
			select {
			case ev := <-eventStream:
				json, err := json.Marshal(ev)
				if err != nil {
					fmt.Println(err)
					continue
				}
				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					fmt.Println(err)
					continue
				}
				w.Write(json)
				err = w.Close()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case <-end:
				go func() {
					for ok := true; ok; _, ok = <-eventStream {
					}
				}()
				game.RemoveListener(listener)
				close(eventStream)
				return
			}
		}
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("error: %v\n", err)
				}
				break
			}
			var baseData ActionBase
			err = json.Unmarshal(message, &baseData)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch baseData.ActionType {
			case ActionTypeHint:
				game.DoHint(playerIdx)
			case ActionTypeDiscard:
				var discardData ActionDiscard
				err = json.Unmarshal(message, &discardData)
				if err != nil {
					fmt.Println(err)
					continue
				}
				game.DoDiscard(playerIdx, discardData.Index)
			case ActionTypePlay:
				var playData ActionPlay
				err = json.Unmarshal(message, &playData)
				if err != nil {
					fmt.Println(err)
					continue
				}
				game.DoPlay(playerIdx, playData.Index)
			case ActionTypeTilt:
				var tiltData ActionTilt
				err = json.Unmarshal(message, &tiltData)
				if err != nil {
					fmt.Println(err)
					continue
				}
				game.DoTilt(playerIdx, tiltData.Index, tiltData.Tilt)
			case ActionTypeMove:
				var moveData ActionMove
				err = json.Unmarshal(message, &moveData)
				if err != nil {
					fmt.Println(err)
					continue
				}
				game.DoMove(playerIdx, moveData.From, moveData.To)
			}
		}
		end <- struct{}{}
	}()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game = NewGame(nplayers)

	http.Handle("/", http.FileServer(http.Dir("./hanabi-frontend/build")))
	http.HandleFunc("/join", handleJoin)
	http.HandleFunc("/connect", handleConnect)

	panic(http.ListenAndServe(":80", nil))
}
