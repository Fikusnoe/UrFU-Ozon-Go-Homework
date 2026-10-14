package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadConsole() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
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
		default:
			fmt.Println("Неизвестная команда!")
		}

		fmt.Println()
		fmt.Print("Введите команду: ")
	}

	err := scanner.Err()
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
	}
}

func PrintInterface() {
	fmt.Println("TaskTracker ver 0.1")
	fmt.Println("Привет!")
	fmt.Println("Команды:")
	fmt.Println("help - полный список команд и их описание")
	fmt.Println("exit - выход из приложения")
	fmt.Print("Введите команду: ")
	ReadConsole()
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
