package config

const apiAddr = "http://localhost:8080/api"
const execAddr = "http://localhost:8090/exec"

func GetAPIAddr() string {
	return apiAddr
}

func GetExecAddr() string {
	return execAddr
}
