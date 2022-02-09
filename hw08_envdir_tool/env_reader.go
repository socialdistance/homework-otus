package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrorDir  = errors.New("error directory")
	ErrNotDir = errors.New("not directory")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, ErrorDir
	}
	if !stat.IsDir() {
		return nil, ErrNotDir
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	envs := Environment{}

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue
		}

		result := fmt.Sprintf("%s%c%s", strings.TrimRight(dir, string(os.PathSeparator)), os.PathSeparator, file.Name())

		statFile, err := os.Stat(result)
		if err != nil {
			return nil, err
		}
		if statFile.IsDir() {
			continue
		}

		if statFile.Size() == 0 {
			envs[file.Name()] = EnvValue{
				NeedRemove: true,
			}
			continue
		}

		firstLine, err := checkFirstLine(result)
		if err != nil {
			return nil, err
		}
		envs[file.Name()] = EnvValue{
			Value: clear(firstLine),
		}
	}

	return envs, nil
}

func checkFirstLine(file string) (string, error) {
	openFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer openFile.Close()

	data, err := ioutil.ReadAll(openFile)
	if err != nil {
		return "", err
	}

	res := strings.Split(string(data), "\n")[0]

	return res, nil
}

func clear(str string) string {
	rep := strings.ReplaceAll(str, "\x00", "\n")
	res := strings.TrimRight(rep, "\t ")

	return res
}
