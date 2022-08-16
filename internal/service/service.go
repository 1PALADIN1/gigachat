package service

var signingKey = ""

func Init(key string) {
	signingKey = key
}
