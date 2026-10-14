package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"hw1/internal/app"
)

func ReadConsole(taskHandler app.TaskHandler) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("Введите команду: ")
		fmt.Println()

		rawQueryStr := scanner.Text()
		querySlice := strings.Fields(rawQueryStr)

		command := querySlice[0]
		switch command {
		case "help":
			PrintHelp()
		case "exit":
			return
		case "statuses":
			PrintStatuses()
		case "clear":
			fmt.Print("\033[2J")
			PrintInterface()
		case "add":
			taskHandler.AddTask()
		case "status":
			taskHandler.UpdateTaskStatus()
		case "edit":
			taskHandler.EditTask()
		case "delete":
			taskHandler.DeleteTask()
		default:
			fmt.Println("Неизвестная команда!")
		}
	}

	err := scanner.Err()
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
	}
}

func Run(taskHandler app.TaskHandler) {
	PrintInterface()
	ReadConsole(taskHandler)
}

func PrintInterface() {
	fmt.Println("TaskTracker ver 0.1")
	fmt.Println("Привет!")
	fmt.Println("Команды:")
	fmt.Println("help - полный список команд и их описание")
	fmt.Println("exit - выход из приложения")
	fmt.Print("Введите команду: ")
}

func PrintHelp() {
	fmt.Println()
	fmt.Println("Список всех команд:")
	fmt.Println("add <заголовок> <дедлайн> - добавить задачу.")
	fmt.Println("list - вывести все задачи на экран")
	fmt.Println("show <идентификатор_задачи> - вывести на экран задачу.")
	fmt.Println("status <идентификатор_задачи> <новый_статус> - с")
	fmt.Println("statuses - вывести список всех статусов задачи")
	fmt.Println("edit <идентификатор_задачи> <заголовок> <дедлайн> - отредактировать задачу. Можно ввести заголовок и/или дедлайн")
	fmt.Println("delete <идентификатор_задачи> - удалить задачу.")
	fmt.Println("clear - очищает экран, отрисовывает стартовый интерфейс")
	fmt.Println()
	fmt.Println("!!!ВАЖНО!!!")
	fmt.Println("<> писать не нужно")
	fmt.Println("Заголовок - строка, идентификатор задачи - число")
	fmt.Println("Дедлайн записывается в формате 2006-01-02")
	fmt.Println("Пример команды: add пройти тест 2026-10-16")
}

func PrintStatuses() {
	fmt.Println("Список статусов:")
	fmt.Println("planned")
	fmt.Println("in_progress")
	fmt.Println("cancelled")
	fmt.Println("done")
}
