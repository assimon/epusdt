package dao

import (
	"log"
)

func Init() {
	if err := DBInit(); err != nil {
		log.Fatalf("DBInit err: %v", err)
	}

	if err := RedisInit(); err != nil {
		log.Fatalf("RedisInit err: %v", err)
	}
}
