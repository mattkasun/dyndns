package main

import (
	"encoding/json"
	"fmt"
)

type record struct {
  Domain_Record struct {
    Name string
    Id int
    Data string
  }
}

func main() {
  //b := []byte(`{"id": 47530833,"name": "winterfell", "data": "45.72.157.43"}`)
  b := []byte(`{"domain_record": {"id":47530833,"type":"A","name":"winterfell","data":"45.72.157.43","priority":"null","port":"null","ttl":3600,"weight":"null","flags":"null","tag":"null"}}`)

	var dat record

	if err := json.Unmarshal(b, &dat); err != nil {
		panic(err)
	}
  //n := record{}
  fmt.Println(dat)
}
