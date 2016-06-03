package util

import "os"

const DefaultPort = "3000"

func GetPort() string {
	var port = os.Getenv("PORT")
	if port != "" {
		return port
	}
	return DefaultPort
}
