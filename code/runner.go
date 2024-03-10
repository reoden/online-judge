package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	//go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(pipe, "23 11\n")
	if err := cmd.Run(); err != nil {
		log.Fatalln(err, stderr.String())
	}
	fmt.Println(out.String())
	println(out.String() == "34\n")
}
