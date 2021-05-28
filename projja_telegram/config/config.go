package config

const apiAddr = "http://backend-api:8080/api"
const execAddr = "http://backend-exec:8090/exec"

func GetAPIAddr() string {
	return apiAddr
}

func GetExecAddr() string {
	return execAddr
}
