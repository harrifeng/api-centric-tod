package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func multiReplace(inFile string, outFile string, replace map[string]string) bool {

	return false
}

func readlines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLinesAndReplace(lines []string, path string, reps map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		for k, v := range reps {
			line = strings.Replace(line, k, v, -1)
		}
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func createNewConfigFile(srcPath string, newPath string, reps map[string]string) error {
	lines, err := readlines(srcPath)

	if err != nil {
		return err
	}

	return writeLinesAndReplace(lines, newPath, reps)
}

func main() {

	var src string
	var port string
	var commit string
	var servername string

	flag.StringVar(&src, "src", "ct", "source folder location")
	flag.StringVar(&port, "port", "13000", "default port number")
	flag.StringVar(&commit, "commit", "HEAD", "default git commit ID")

	flag.StringVar(&servername, "servername", "localhost", "server dns name")
	flag.Parse()

	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	dockerFilePath := filepath.Join(rootPath, "Dockerfile.tmp")
	dockerComposePath := filepath.Join(rootPath, "docker-compose.yml.tmp")

	newDockerFilePath := "Dockerfile-" + port + ".yml"
	newDockerComposePath := "docker-compose-" + port + ".yml"

	createNewConfigFile(dockerFilePath, newDockerFilePath, map[string]string{"RTM_SERVER_NAME_PORT": servername + ":" + port, "RTM_SOURCE_FOLDER": src, "RTM_COMMIT_ID": commit, "RTM_DB_HOST": port + "_db_1"})
	createNewConfigFile(dockerComposePath, newDockerComposePath, map[string]string{"RTM_PORT": port, "RTM_DOFILE": newDockerFilePath, "RTM_DB_HOST": port + "_db_1"})

	cmd := exec.Command("docker-compose", "-f", newDockerComposePath, "-p", port, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
