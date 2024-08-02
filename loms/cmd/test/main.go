package main

import (
	"github.com/google/uuid"
	"log"
)

func main() {

	u := uuid.New()
	log.Println(u.String())

	var b []byte

	b = u[:]
	log.Println(b)

	u1, err := uuid.FromBytes(b)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(u1.String())

}
