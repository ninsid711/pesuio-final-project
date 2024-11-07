package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/anuragrao04/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	var request models.RunRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "wrong request"})
		return
	}

	var output string
	var err error
	switch request.Language {
	case "python":
		output, err = execPythonCode(request.Code, request.Input)
	case "go":
		output, err = execGolangCode(request.Code, request.Input)
	case "C":
		output, err = execClangCode(request.Code, request.Input)
	default:
		c.JSON(400, gin.H{"error": "unsupported language"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"output": output})
}

func execPythonCode(code, input string) (string, error) {
	cmd := exec.Command("python", "-c", code)
	var stdin bytes.Buffer
	stdin.WriteString(input)
	cmd.Stdin = &stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s\nError: %v", out.String(), err)
	}
	return out.String(), nil
}

func execClangCode(code, input string) (string, error) {
	// Write C code to a temporary file
	tmpFile, err := os.CreateTemp("", "code-*.c")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.WriteString(tmpFile, code); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Compile the C code
	outputFile := tmpFile.Name() + ".out"
	cmd := exec.Command("gcc", tmpFile.Name(), "-o", outputFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s\nError: %v", out, err)
	}
	defer os.Remove(outputFile)

	// Run the compiled executable with input
	cmd = exec.Command(outputFile)
	var stdin bytes.Buffer
	stdin.WriteString(input)
	cmd.Stdin = &stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s\nError: %v", out.String(), err)
	}
	return out.String(), nil
}

// Function to execute Go code
func execGolangCode(code, input string) (string, error) {
	// Write Go code to a temporary file
	tmpFile, err := os.CreateTemp("", "code-*.go")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.WriteString(tmpFile, code); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Run the Go code using `go run`
	cmd := exec.Command("go", "run", tmpFile.Name())
	var stdin bytes.Buffer
	stdin.WriteString(input)
	cmd.Stdin = &stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s\nError: %v", out.String(), err)
	}
	return out.String(), nil
}
