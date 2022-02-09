package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// проверка наличия команд и аргументов
	if len(cmd) == 0 {
		return 0
	}

	// парсим cmd
	command := cmd[0]
	args := cmd[1:]

	// получаем структуру для выполнения команды
	com := exec.Command(command, args...)
	// пробрасываем потоки ввода/вывода и ошибок
	com.Stdout = os.Stdout
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin

	// берем окружение
	comEnv := os.Environ()
	for fileName, v := range env {
		// сверяем параметры
		for i, str := range comEnv {
			// если параметр совпадает с именем файла
			if strings.Split(str, "=")[0] == fileName {
				// удаляем
				comEnv = append(comEnv[:i], comEnv[i+1:]...)
			}
		}
		if !v.NeedRemove {
			// если не помечен на удаление, то добавляем в конец
			comEnv = append(comEnv, fmt.Sprintf("%s=%s", fileName, v.Value))
		}
	}
	// присваиваем окружение команде
	com.Env = comEnv

	// запускаем
	if err := com.Run(); err != nil {
		return 1
	}

	return 0
}
