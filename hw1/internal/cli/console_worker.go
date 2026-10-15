package cli

import (
	"bufio"
	"errors"
	"fmt"
	"hw1/internal/app"
	"hw1/internal/model"
	"os"
	"strings"
	"time"
)

var (
	ErrUnknownCommand       = errors.New("Неизвестная команда")
	ErrInvalidCommandFormat = errors.New("Неверный формат команды")
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
			err := AddTask(taskHandler, querySlice, querySliceLen)
			if err != nil {
				PrintError(err)
			}

		case "list":
			err := ListTasks(taskHandler, querySlice)
			if err != nil {
				PrintError(err)
			}

		case "show":
			err := ShowTaskById(taskHandler, querySlice, querySliceLen)
			if err != nil {
				PrintError(err)
			}

		case "status":
			err := ChangeTaskStatus(taskHandler, querySlice, querySliceLen)
			if err != nil {
				PrintError(err)
			}

		case "edit":
			err := EditTask(taskHandler, querySlice, querySliceLen)
			if err != nil {
				PrintError(err)
			}

		case "delete":
			err := DeleteTask(taskHandler, querySlice, querySliceLen)
			if err != nil {
				PrintError(err)
			}

		default:
			fmt.Println(ErrUnknownCommand.Error())
		}
		fmt.Println()
		fmt.Print("Введите команду: ")
	}

	err := scanner.Err()
	if err != nil {
		PrintError(err)
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

func PrintError(err error) {
	fmt.Println("Ошибка: " + err.Error())
}

func ValidateDeadline(deadline time.Time) error {
	if deadline.After(time.Now()) {
		return fmt.Errorf("Неверный формат даты: дедлайн не может быть позже текущей даты")
	}
	return nil
}

func AddTask(taskHandler app.TaskHandler, querySlice []string, querySliceLen int) error {
	if querySliceLen < 3 {
		return ErrInvalidCommandFormat
	}

	name := strings.Join(querySlice[1:querySliceLen-2], " ")
	deadline, err := time.Parse("2006-01-02", querySlice[querySliceLen-1])
	if err != nil {
		return err
	}

	err = ValidateDeadline(deadline)
	if err != nil {
		return err
	}

	err = taskHandler.AddTask(name, deadline)
	if err != nil {
		return err
	}

	return nil
}

func ListTasks(taskHandler app.TaskHandler, querySlice []string) error {
	tasks, err := taskHandler.GetAllTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Println(task.String())
	}

	return nil
}

func ShowTaskById(taskHandler app.TaskHandler, querySlice []string, querySliceLen int) error {
	if querySliceLen < 2 || querySliceLen > 2 {
		return ErrInvalidCommandFormat
	}

	taskId, err := model.ParseID(querySlice[1])
	if err != nil {
		return err
	}

	task, err := taskHandler.GetTaskById(taskId)
	if err != nil {
		return err
	}

	fmt.Println(task.String())

	return nil
}

func ChangeTaskStatus(taskHandler app.TaskHandler, querySlice []string, querySliceLen int) error {
	if querySliceLen < 2 {
		return ErrInvalidCommandFormat
	}

	taskId, err := model.ParseID(querySlice[1])
	if err != nil {
		return err
	}

	status, err := model.ParseStatus(querySlice[2])
	if err != nil {
		return err
	}

	err = taskHandler.UpdateTaskStatus(taskId, status)
	if err != nil {
		return err
	}

	return nil
}

func EditTask(taskHandler app.TaskHandler, querySlice []string, querySliceLen int) error {
	if querySliceLen < 3 {
		return ErrInvalidCommandFormat
	}

	taskId, err := model.ParseID(querySlice[1])
	if err != nil {
		return err
	}

	deadline, err := time.Parse("2006-01-02", querySlice[querySliceLen-1])
	if err != nil {
		name := strings.Join(querySlice[1:], " ")
		err = taskHandler.EditTask(taskId, name, time.Time{})
		if err != nil {
			return err
		}
	}

	name := ""
	if querySliceLen > 4 {
		name = strings.Join(querySlice[1:querySliceLen-2], " ")
	}

	err = taskHandler.EditTask(taskId, name, deadline)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(taskHandler app.TaskHandler, querySlice []string, querySliceLen int) error {
	if querySliceLen < 2 {
		return ErrInvalidCommandFormat
	}

	taskId, err := model.ParseID(querySlice[1])
	if err != nil {
		return err
	}

	err = taskHandler.DeleteTask(taskId)
	if err != nil {
		return err
	}

	return nil
}
