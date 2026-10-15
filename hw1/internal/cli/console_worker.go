package cli

import (
	"bufio"
	"fmt"
	"hw1/internal/app"
	"hw1/internal/model"
	"os"
	"strings"
	"time"
)

func ReadConsole(taskHandler app.TaskHandler) {
	
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println()
	fmt.Print("Введите команду: ")
	
	for scanner.Scan() {
		rawQueryStr := scanner.Text()
		querySlice := strings.Fields(rawQueryStr)
		querySliceLen := len(querySlice)

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
			if querySliceLen < 3 {
				fmt.Println("Неверный формат команды")
			}

			name := strings.Join(querySlice[1:querySliceLen-2], " ")
			deadline, _ := time.Parse("2006-01-02", querySlice[querySliceLen-1])
			if deadline.After(time.Now()) {
				fmt.Println("Неверный формат даты: дедлайн не может быть позже текущей даты")
			}

			taskHandler.AddTask(name, deadline) // обработать ошибку

		case "list":
			tasks, _ := taskHandler.GetAllTasks() // обработать ошибку
			for _, task := range tasks {
				fmt.Println(task.String())
			}

		case "show":
			if querySliceLen < 2 || querySliceLen > 2 {
				fmt.Println("Неверный формат команды")
				break
			}

			taskId, _ := model.ParseID(querySlice[1])
			task, _ := taskHandler.GetTaskById(taskId) // обработать ошибку
			fmt.Println(task.String())

		case "status":
			if querySliceLen < 2 {
				fmt.Println("Неверный формат команды")
				break
			}

			taskId, _ := model.ParseID(querySlice[1])
			status, _ := model.ParseStatus(querySlice[2])
			err := taskHandler.UpdateTaskStatus(taskId, status)
			if err != nil {
				fmt.Println(err.Error())
			}

		case "edit":
			if querySliceLen < 3 {
				fmt.Println("Неверный формат команды")
			}

			taskId, _ := model.ParseID(querySlice[1])

			deadline, errDeadline := time.Parse("2006-01-02", querySlice[querySliceLen-1])
			if errDeadline != nil {
				name := strings.Join(querySlice[1:], " ")
				taskHandler.EditTask(taskId, name, time.Time{}) // обработать ошибку
			} 
			
			name := ""
			if querySliceLen > 4 {
				name = strings.Join(querySlice[1:querySliceLen-2], " ")
			} 
			
			taskHandler.EditTask(taskId, name, deadline) // обработать ошибку

		case "delete":
			if querySliceLen < 2 {
				fmt.Println("Неверный формат команды")
			}

			taskId, _ := model.ParseID(querySlice[1])
			taskHandler.DeleteTask(taskId) // обработать ошибку

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

func Run(taskHandler app.TaskHandler) {
	PrintInterface()
	ReadConsole(taskHandler)
}

func PrintInterface() {
	fmt.Println("TaskTracker ver 0.5")
	fmt.Println("Привет!")
	fmt.Println("Команды:")
	fmt.Println("help - полный список команд и их описание")
	fmt.Println("exit - выход из приложения")
}

func PrintHelp() {
	fmt.Println()
	fmt.Println("Список всех команд:")
	fmt.Println("add <заголовок> <дедлайн> - добавить задачу.")
	fmt.Println("list - вывести все задачи на экран")
	fmt.Println("show <идентификатор_задачи> - вывести на экран задачу.")
	fmt.Println("status <идентификатор_задачи> <новый_статус> - с")
	fmt.Println("statuses - вывести список всех возможных статусов задачи")
	fmt.Println("edit <идентификатор_задачи> <заголовок> <дедлайн> - отредактировать задачу. Можно ввести заголовок и/или дедлайн")
	fmt.Println("delete <идентификатор_задачи> - удалить задачу.")
	fmt.Println("clear - очищает экран, отрисовывает стартовый интерфейс")
	fmt.Println()
	fmt.Println("!!!ВАЖНО!!!")
	fmt.Println("<> писать не нужно")
	fmt.Println("Заголовок - строка, идентификатор задачи - число")
	fmt.Println("Если дата будет введена неверно - она засчитается как текст")
	fmt.Println("Пример команды: add пройти тест 2026-10-16")
}

func PrintStatuses() {
	fmt.Println("Список статусов:")
	fmt.Println("planned")
	fmt.Println("in_progress")
	fmt.Println("cancelled")
	fmt.Println("done")
}
