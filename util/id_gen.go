package util

import (
	"os"
	"strconv"
	"time"
)

var IdGenerator = make(chan (chan int64)) // NOTE: only for user id
var sequence uint8
var pid int

func StartUserIdGenerator() {
	pid = os.Getpid()
	sequence = 0
	for req := range IdGenerator {
		timestamp := time.Now().UnixNano()
		newId := timestamp | int64(pid)<<22
		newId |= int64(sequence)
		newIdString := strconv.FormatInt(newId, 10)
		newIdString = newIdString[5:]
		newId, _ = strconv.ParseInt(newIdString, 10, 64)
		sequence += 1
		req <- newId
	}
}

func GetID(generatorChannel chan (chan int64)) int64 {
	waitChan := make(chan int64)
	generatorChannel <- waitChan

	t := time.NewTimer(10 * time.Millisecond)

	select {
	case id := <-waitChan:
		close(waitChan)
		t.Stop()
		return id
	case <-t.C:
		panic("generate new user id too slow")
	}

	// Output:
	// 100000
}