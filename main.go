package main

import (
	"fmt"
	"log"
	"os/exec"
	"bytes"
)


func runCommandWithOutput(command string) (string) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Fatal("An invalid command was executed: ", command)
	}
	return fmt.Sprintf("%s", out)
}

func runffmpeg(command string, params string){
	cmd := exec.Command(command, params)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// Redirigir el stdErr una variable para luego poder imprimirlo.
	
	err := cmd.Run()

	if err != nil{
		log.Fatal("Command finished with error \n", stderr.String())
	}
	log.Printf("Command finished succesfulyy")
}

func main(){
	/*
		var command1 string = "date"
	out  := runCommandWithOutput(command1)
	fmt.Println(out)
	command2 := "dater"
	out2  := runCommandWithOutput(command2)
	fmt.Println(out2)
	*/

	mkdirCommand := "mkdir"
	params := "hello boiler"
	runffmpeg(mkdirCommand, params)

	
}