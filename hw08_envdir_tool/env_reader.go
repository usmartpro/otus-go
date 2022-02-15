package main

import (
	"bytes"
	"errors"
	"os"
	"path"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrUnknownPath = errors.New("unknown path")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// считываем список папок
	listFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrUnknownPath
	}

	// инициализируем мапу для ответа
	newEnvs := make(Environment)
	for _, currentFile := range listFiles {
		// инфо о файле
		currentFileInfo, err := currentFile.Info()
		if err != nil {
			return nil, err
		}

		// пустой файл
		if currentFileInfo.Size() == 0 {
			// помечаем имя файла к удалению
			newEnvs[currentFile.Name()] = EnvValue{"", true}
			continue
		}

		// считываем данные файла
		fileData, err := os.ReadFile(path.Join(dir, currentFile.Name()))
		if err != nil {
			return nil, err
		}

		// чистим от лишних символов
		fileData = bytes.Split(fileData, []byte("\n"))[0]
		fileData = bytes.ReplaceAll(fileData, []byte{0}, []byte("\n"))
		fileData = bytes.TrimRight(fileData, " \t")

		// добавляем к результату
		newEnvs[currentFile.Name()] = EnvValue{string(fileData), false}
	}

	return newEnvs, nil
}
