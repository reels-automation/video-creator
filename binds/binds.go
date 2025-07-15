package binds
import (
	"log"
	"os/exec"
	"bytes"
	"fmt"
)

func RunCommand(command string, params []string){
	cmd := exec.Command(command, params...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// Redirigir el stdErr una variable para luego poder imprimirlo.
	
	err := cmd.Run()

	if err != nil{
		log.Fatal("Command finished with error \nGo Err: \n", err, "\nStdErr: \n", stderr.String())
	}
	log.Printf("Command finished succesfulyy")
}

func RunCommandWithOutput(command string, params []string) (string) {
	out, err := exec.Command(command, params...).Output()
	if err != nil {
		log.Fatal("An invalid command was executed: ", command)
	}
	return fmt.Sprintf("%s", out)
}