package binds
import (
	"log"
	"os/exec"
	"bytes"
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

func RunCommandWithOutput(command string, params []string) string {
	cmd := exec.Command(command, params...)
	out, err := cmd.CombinedOutput() // captures both stdout and stderr

	if err != nil {
		log.Fatalf(
			"Command execution failed:\nCommand: %s\nArgs: %v\nError: %v\nOutput: %s",
			command, params, err, string(out),
		)
	}

	return string(out)
}
