package main

import (
	"fmt"
	"strings"
)

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные перменные ( тут можно ), функции
*/

type Location struct {
	name   string
	moving []string
	items  []string
}

var (
	gameMap  map[string]Location
	location Location
	items    []string
	door     bool
	backpack bool
	subject  []string
)

func useItem(cmd []string) string {
	if !find(items, cmd[1]) {
		return fmt.Sprintf("нет предмета в инвентаре - %s", cmd[1])
	}
	if !find(subject, cmd[2]) {
		return "не к чему применить"
	}
	if location.name == "коридор" && cmd[1] == "ключи" && cmd[2] == "дверь" {
		door = true
		return "дверь открыта"
	}

	return fmt.Sprintf("вы применили %s к %s", cmd[1], cmd[2])
}

func takeItem(cmd []string) string {
	if !find(location.items, cmd[1]) {
		return "нет такого"
	}
	if !backpack {
		return "некуда класть"
	}
	items = append(items, cmd[1])
	for ind := range location.items {
		if location.items[ind] == cmd[1] {
			location.items = remove(location.items, ind)
			break
		}
	}
	return fmt.Sprintf("предмет добавлен в инвентарь: %s", cmd[1])
}

func putBackpack(cmd []string) string {
	if cmd[1] != "рюкзак" {
		return fmt.Sprintf("%s нельзя надеть", cmd[1])
	}
	if backpack {
		return "рюкзак уже надет"
	}
	if location.name != "комната" {
		return "рюкзака тут нет"
	}
	backpack = true
	location.items = remove(location.items, 2)
	return "вы надели: рюкзак"
}

func look() string {
	switch location.name {
	case "кухня":
		if backpack {
			return fmt.Sprintf("ты находишься на кухне, на столе: %s, надо идти в универ. можно пройти - %s",
				toString(location.items), toString(location.moving))
		}
		return fmt.Sprintf("ты находишься на кухне, на столе: %s, надо собрать рюкзак и идти в универ. можно пройти - %s",
			toString(location.items), toString(location.moving))
	case "комната":
		if len(location.items) == 0 {
			return fmt.Sprintf("пустая комната. можно пройти - %s",
				toString(location.moving))
		}
		if backpack {
			return fmt.Sprintf("на столе: %s. можно пройти - %s",
				toString(location.items), toString(location.moving))
		}
		return fmt.Sprintf("на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - %s",
			toString(location.moving))
	default:
		return fmt.Sprintf("%s, ничего интересного. можно пройти - %s", location.name, toString(location.moving))
	}
}

func going(command []string) string {
	if find(location.moving, command[1]) {
		location = gameMap[command[1]]
		switch location.name {
		case "кухня":
			return fmt.Sprintf("%s, ничего интересного. можно пройти - %s", location.name, toString(location.moving))
		case "коридор":
			return fmt.Sprintf("ничего интересного. можно пройти - %s", toString(location.moving))
		case "комната":
			return fmt.Sprintf("ты в своей комнате. можно пройти - %s", toString(location.moving))
		case "улица":
			if door {
				return fmt.Sprintf("на улице весна. можно пройти - %s", toString(location.moving))
			}
			location = gameMap["коридор"]
			return "дверь закрыта"
		}
	}
	return fmt.Sprintf("нет пути в %s", command[1])
}

func toString(strings []string) string {
	str := ""
	for ind := range strings {
		if ind == len(strings)-1 {
			str += strings[ind]
			break
		}
		str += strings[ind] + ", "
	}
	return str
}

func remove[T any](slice []T, i int) []T {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func find(array []string, element string) bool {
	for ind := range array {
		if array[ind] == element {
			return true
		}
	}
	return false
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
}

func initGame() {
	/*
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/
	gameMap = map[string]Location{
		"кухня": Location{
			"кухня", []string{"коридор"}, []string{"чай"},
		},
		"коридор": Location{
			"коридор",
			[]string{
				"кухня",
				"комната",
				"улица",
			},
			[]string{},
		},
		"комната": Location{
			"комната", []string{"коридор"}, []string{"ключи", "конспекты", "рюкзак"},
		},
		"улица": Location{
			"улица", []string{"домой"}, []string{},
		},
	}
	location = gameMap["кухня"]
	door = false
	backpack = false
	subject = []string{"дверь"}
	items = make([]string, 0, 2)
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/

	cmd := strings.Split(command, " ")
	if cmd[0] == "идти" {
		return going(cmd)
	}
	if cmd[0] == "осмотреться" {
		return look()
	}
	if cmd[0] == "надеть" {
		return putBackpack(cmd)
	}
	if cmd[0] == "взять" {
		return takeItem(cmd)
	}
	if cmd[0] == "применить" {
		return useItem(cmd)
	}

	return "неизвестная команда"
}
