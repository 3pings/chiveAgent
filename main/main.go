package main

import (
	"fmt"
	"log"
	"os"

	"github.com/udhos/acigo/aci"
)

func main() {

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
		//name := n["name"]
		dn := n["dn"]
		role := n["role"]
		serial := n["serial"]
		fmt.Println(dn, role, serial)
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
