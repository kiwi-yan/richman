package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"syscall"
	"time"
)

type Assets struct{}

func getUser(name string) *User {
	return userMap[name]
}

var (
	userMap     map[string]*User
	cmdMessage  chan *Command
	messageIn   chan string
	userMessage chan *Command
	lands       []*Land
)

type Land struct {
	owner     string
	used      bool
	seed      Object
	startTime time.Time
}

func main() {
	fakeUser()
	initResource()
	initUserModule()

	initMessageModule()
	initInteractModule()
	runTheWorld()
	// main2()
}

func runTheWorld() {
	for {
		time.Sleep(1 * time.Second)
		printStatInfo()
	}
}

func printStatInfo() {
}

func initResource() {
	lands = make([]*Land, 64)
}

func initInteractModule() {
	messageIn = make(chan string)
	cmdMessage = make(chan *Command)
	go func() {
		for {
			var cmd string
			_, err := fmt.Scan(&cmd)
			if err == io.EOF {
				break
			}
			fmt.Println(">", cmd)
			c, err := ParseCommand(cmd)
			if err != nil {
				fmt.Println("ParseCommand error:", err)
				continue
			}
			cmdMessage <- c
			in := <-messageIn
			fmt.Println(in)
		}
	}()
}

func initMessageModule() {
	userMessage = make(chan *Command)
	go func() {
		for {
			cmd := <-cmdMessage
			switch cmd.Module {
			case "user":
				userMessage <- cmd
			case "game":
				if cmd.Action == "exit" {
					os.Exit(0)
				}
			default:
				messageIn <- fmt.Sprintln("unknown module:", cmd.Module)
			}
		}
	}()
}

func initUserModule() {
	userMap = make(map[string]*User, 0)
	go func() {
		for {
			cmd := <-userMessage
			args := cmd.Args.(*UserCommand)
			switch cmd.Action {
			case "register":
				userMap[args.UserName] = NewUser(args.UserName)
				messageIn <- fmt.Sprintln("register user", args.UserName, "success!")
			case "get":
				messageIn <- fmt.Sprintln(getUser(args.UserName))
			case "exec":
				messageIn <- execUserAction(args)
			default:
				messageIn <- fmt.Sprintln("unknown action:", cmd.Action)
			}
		}
	}()
}

func execUserAction(userCmd *UserCommand) string {
	user := getUser(userCmd.UserName)
	if user == nil {
		return fmt.Sprintln("Cannot find user:", userCmd.UserName)
	}
	switch userCmd.UserAction {
	case "buy":
		if user.Buy(userCmd.Object, userCmd.Number) {
			return fmt.Sprintf("User %s buy %d %s\n", userCmd.UserName, userCmd.Number, reflect.TypeOf(userCmd.Object).Name())
		}
		return fmt.Sprintf("User %s's money is not enough\n", userCmd.UserName)
	case "sell":
		if user.Sell(userCmd.Object, userCmd.Number) {
			return fmt.Sprintf("User %s sell %d %s\n", userCmd.UserName, userCmd.Number, reflect.TypeOf(userCmd.Object).Name())
		}
		return fmt.Sprintf("User %s's %s is not enough\n", userCmd.UserName, reflect.TypeOf(userCmd.Object).Name())
	default:
		return fmt.Sprintln("user:exec unsupported action:", userCmd.UserAction)
	}
}

var ObjectStringMap = map[string]Object{
	"food": Food{},
}

func convertToObject(s string) (obj Object, err error) {
	var ok bool
	obj, ok = ObjectStringMap[s]
	if !ok {
		err = fmt.Errorf("Cannot convert to object: %s", s)
	}
	return
}

func fakeUser() {
	file, _ := os.Open("command.txt")
	syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
}
