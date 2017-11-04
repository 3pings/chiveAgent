package utility

import (
	"bytes"
	"fmt"
	"github.com/vallard/spark"
	"io/ioutil"
	"net/http"
)

//SendJSON allows us to send JSON to a remote device

func SendJSON(jsonStr []byte, url string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// SendSparkMessage to Room
func SendSparkMessage(sparkToken, sparkRoomId, sparkMessage string) {
	s := spark.New(sparkToken)
	m := spark.Message{
		RoomId: sparkRoomId,
		Text:   sparkMessage,
	}
	// Post the message to the room
	_, err := s.CreateMessage(m)

	if err != nil {
		panic(err)
	}

}
