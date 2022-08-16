package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)


var actionList []string
func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	wind := []string{"N", "E", "S", "W"}
	var action string
	me, ok := input.Arena.State[input.Links.Self.Href]
	if ok {
		delete(input.Arena.State, input.Links.Self.Href)
	}
	for _, player_state := range input.Arena.State {
		if player_state.X  == me.X || player_state.Y == me.Y {
			if (player_state.X < me.X && me.Direction == "W") || (player_state.X > me.X && me.Direction == "E") || (player_state.Y < me.Y && me.Direction == "N") || (player_state.Y > me.Y && me.Direction == "S"){
				actionList = append(actionList, "T")
				break
			}
			if (player_state.X < me.X || player_state.X > me.X) && (player_state.Direction == "E" || player_state.Direction == "W") {
				if me.Direction == "E" || me.Direction == "W" {
					commands := []string{"L", "R"}
					rand := rand2.Intn(2)
					actionList = append(actionList, commands[rand])
					actionList = append(actionList, "F")
				}else {
					actionList = append(actionList, "F")
				}
				break
			}else if (player_state.Y < me.Y || player_state.Y > me.Y) && (player_state.Direction == "N" || player_state.Direction == "S") {
				if me.Direction == "N" || me.Direction == "S" {
					commands := []string{"L", "R"}
					rand := rand2.Intn(2)
					actionList = append(actionList, commands[rand])
					actionList = append(actionList, "F")
				}else {
					actionList = append(actionList, "F")
				}
				break
			}
		}
	}
	if len(actionList) == 0 {
		commands := []string{"F", "R", "L", "T"}
		rand := rand2.Intn(4)
		actionList = append(actionList, commands[rand])
	}
	action, actionList = actionList[0], actionList[1:]
	return action
}
