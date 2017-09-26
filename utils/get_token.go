package utils

import "fmt"

func GetToken(apicIP string) {
	baseUrl := "http://" + apicIP + "/api/"
	fmt.Println(baseUrl)
}
