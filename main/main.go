package main

import (
	"fmt"
	"github.com/3pings/acigo/aci"
	"log"
	"os"
	"time"
)

func main() {

	type nodeData struct {
		temp string
	}

	nodeInfo := make(map[string]nodeData)

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
				nodeInfo[name] = nodeData{tempMax}
			}
		}
	}
	//Printing today need to add api call
	fmt.Println(nodeInfo)

	// Loop through infinitely every X number of seconds
	for {
		time.Sleep(60 * time.Second)
		errRefresh := a.Refresh()
		if errRefresh != nil {
			log.Printf("refresh %d/%d error: %v", errRefresh)
			os.Exit(3)
		}

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
					nodeInfo[name] = nodeData{tempMax}
				}
			}
		}
	}
	//Printing today need to add api call
	fmt.Println(nodeInfo)
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
