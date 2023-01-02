package main

import "fmt"

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	fmt.Println("Hello4")
}
