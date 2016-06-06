package util

import "os"

const DefaultPort = "3000"
const DefaultUser = "user"
const DefaultPass = "pass"

func GetPort() string {
	var port = os.Getenv("PORT")
	if port != "" {
		return port
	}
	return DefaultPort
}

func GetUser() string {
	var user = os.Getenv("SECURITY_USER_NAME")
	if user != "" {
		return user
	}
	return DefaultUser
}

func GetPassword() string {
	var password = os.Getenv("SECURITY_USER_PASSWORD")
	if password != "" {
		return password
	}
	return DefaultPass
}
