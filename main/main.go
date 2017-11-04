package main

import (
	"encoding/json"
	"fmt"
	"github.com/3pings/acigo/aci"
	"github.com/3pings/chiveAgent/utility"
	"log"
	"os"
	"time"
)

var nodeInfo = make(map[string]string)

func main() {

	token := os.Getenv("SPARKTOKEN")
	roomID := os.Getenv("SPARKROOMID")

	// Get environment variables for APIC login
	debug := os.Getenv("DEBUG") != ""

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	// display existing nodes
	nodes, errList := a.NodeList()
	if errList != nil {
		log.Printf("could not list nodes: %v", errList)
		return
	}

	// loop through to get temperature data per node
	for {

		for _, n := range nodes {

			nodeDetails, errList := a.GetTemp(n["dn"].(string))
			if errList != nil {
				log.Printf("could not list node details: %v", errList)
				return
			}

			for _, d := range nodeDetails {
				name := n["name"].(string)
				tempMax := d["currentMax"].(string)
				nodeRole := n["role"].(string)

				// At this time we do not want the controller info
				if nodeRole != "controller" {
					nodeInfo[name] = tempMax
				}
			}

		}
		//Put results of node data collection into json
		//Printing today need to add api call
		jsonNode, _ := json.Marshal(nodeInfo)
		utility.SendSparkMessage(token, roomID, string(jsonNode))

		// wait a defined number of seconds before looping back through
		time.Sleep(1 * time.Second)
		errRefresh := a.Refresh()
		if errRefresh != nil {
			log.Println(errRefresh)
			os.Exit(3)
		}
	}

}

func login(debug bool) (*aci.Client, error) {

	a, errNew := aci.New(aci.ClientOptions{Debug: false})
	if errNew != nil {
		return nil, fmt.Errorf("login new client error: %v", errNew)
	}

	// Since credentials have not been specified explicitly under ClientOptions,
	// Login() will use env vars: APIC_HOSTS=host, APIC_USER=username, APIC_PASS=pwd
	errLogin := a.Login()
	if errLogin != nil {
		return nil, fmt.Errorf("login error: %v", errLogin)
	}

	return a, nil
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}

	log.Printf("logout: done")
}
