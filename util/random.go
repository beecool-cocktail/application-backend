package util

import(
	"math/rand"
	"time"
)



var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	Init()
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNumberRunes(n int) string {
	Init()
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}

func GenString(offset int) string {
	return RandStringRunes(offset)
}

func GenNumber(offset int) string {
	return RandNumberRunes(offset)
}