package hello_world

import "fmt"

func IncomingData(data string) (string, error) {
	fmt.Println("Hello world! from a plugin")

	return "", nil
}