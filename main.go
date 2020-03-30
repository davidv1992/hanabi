package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// GameRoom bla
type GameRoom struct {
	sync.Mutex
	game       *GameState
	playersIn  int
	nplayers   int
	lastActive time.Time
}

// PlayerTokenData bla
type PlayerTokenData struct {
	playerIdx int
	room      string
}

var rooms = map[string]*GameRoom{}
var playerTokens = map[string]PlayerTokenData{}
var globalLock sync.Mutex

func doCleanup() {
	globalLock.Lock()
	defer globalLock.Unlock()
	removed := map[string]struct{}{}
	for k, v := range rooms {
		if len(v.game.listeners) == 0 && time.Since(v.lastActive) > 1*time.Hour {
			removed[k] = struct{}{}
			delete(rooms, k)
		}
	}
	for k, v := range playerTokens {
		if _, ok := removed[v.room]; ok {
			delete(playerTokens, k)
		}
	}
}

func handleJoinForm(w http.ResponseWriter, r *http.Request) {
	room := r.URL.RawQuery
	if room == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(`
	<html>
	<head>
	<title>Join Hanabi</title>
	</head>
	<body>
	<form method="POST" action="/hanabi/dojoin">
	<label for="name">Name:</label><input type="text" name="name" id="name"/><br/>
	<input type="submit" value="Join"/>
	<input type="hidden" name="room" id="room" value="` + room + `"/>
	</form>
	<p> To invite others, use this link: <a href="https://walkintrack.nl/hanabi/join.html?` + room + `">https://walkintrack.nl/hanabi/join.html?` + room + `</a> </p>
	</body>
	</html>
	`))
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	room := r.FormValue("room")

	if name == "" || room == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	globalLock.Lock()
	rm, ok := rooms[room]
	globalLock.Unlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rm.Lock()
	defer rm.Unlock()

	if rm.playersIn >= rm.nplayers {
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

	rm.game.SetName(rm.playersIn, name)
	rm.playersIn = rm.playersIn + 1

	globalLock.Lock()
	playerTokens[token] = PlayerTokenData{playerIdx: rm.playersIn - 1, room: room}
	globalLock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 24 * 3600,
		Path:   "/",
	})
	w.Header().Add("Location", "/hanabi/play.html")
	w.WriteHeader(http.StatusSeeOther)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	players := r.FormValue("players")
	nplayers, err := strconv.Atoi(players)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if nplayers < 2 || nplayers > 5 {
		fmt.Println(err)
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

	globalLock.Lock()
	rooms[token] = &GameRoom{
		game:       NewGame(nplayers),
		nplayers:   nplayers,
		playersIn:  0,
		lastActive: time.Now(),
	}
	globalLock.Unlock()

	w.Header().Add("Location", "/hanabi/join.html?"+token)
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

	globalLock.Lock()
	playerData, ok := playerTokens[token]
	globalLock.Unlock()

	if !ok {
		conn.Close()
		return
	}

	globalLock.Lock()
	rm := rooms[playerData.room]
	globalLock.Unlock()

	playerIdx := playerData.playerIdx
	game := rm.game

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
				rm.lastActive = time.Now()
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

	go func() {
		for {
			doCleanup()
			time.Sleep(5 * time.Minute)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/join.html", handleJoinForm)
	http.HandleFunc("/create", handleCreate)
	http.HandleFunc("/dojoin", handleJoin)
	http.HandleFunc("/connect", handleConnect)

	panic(http.ListenAndServe(":8081", nil))
}
