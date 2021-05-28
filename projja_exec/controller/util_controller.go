package controller

import (
	"encoding/json"
	"log"
)

type response struct {
	Description string
	Content     interface{}
}

func (c *controller) makeContentResponse(code int, desc string, content interface{}) (int, string) {
	response := &response{
		desc,
		content,
	}
	byteResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error during content marshalling:", err.Error())
		return 500, err.Error()
	}
	return code, string(byteResponse)
}
