package client

import (
	"bufio"
	"fmt"
	"os"
)

// Чтение имени пользователя
func ReadUsername() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите имя пользователя: ")
	scanner.Scan()
	return scanner.Text()
}

// Чтение действия пользователя
func ReadUserAction() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nВыберите действие:\n1 - Отправить сообщение\n2 - Лайкнуть сообщение\n3 - Комментировать сообщение\n\n5 - Выход\n")
	scanner.Scan()
	return scanner.Text()
}
