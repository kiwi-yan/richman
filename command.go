package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	Module string
	Action string
	Args   interface{}
}

type UserCommand struct {
	UserName   string
	UserAction string
	Object     Object
	Number     int
}

var Module2Command = map[string]interface{}{
	"user": UserCommand{},
	"game": struct{}{},
}

var Str2Field = map[string]string{
	"user-name":   "UserName",
	"user-action": "UserAction",
	"user-object": "Object",
	"user-number": "Number",
}

func (c *Command) String() string {
	return fmt.Sprintf(`Module: %s
Action: %s
Args: %#v`, c.Module, c.Action, c.Args)
}

func ParseCommand(cmd string) (*Command, error) {
	cmdList := strings.Split(cmd, ":")
	if len(cmdList) < 2 {
		return nil, errors.New(fmt.Sprint("Error cmd:", cmd))
	}
	c := &Command{}
	c.Module = strings.ToLower(cmdList[0])
	subCmd, ok := Module2Command[c.Module]
	if !ok {
		return nil, errors.New(fmt.Sprint("Unknow module:", cmd))
	}
	c.Action = cmdList[1]
	var err error
	c.Args, err = makeArgs(c.Module, subCmd, cmdList[2:]...)
	return c, err
}

func makeArgs(modlue string, subCmd interface{}, args ...string) (interface{}, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, nil
	}
	// 这里一定得是结构体传进来而非它的指针
	typeOfSubCmd := reflect.TypeOf(subCmd)
	if argsLen%2 != 0 {
		return nil, fmt.Errorf("Subcommand field and value are not matched: %v", args)
	}
	// 使用类型创建一个对应的指针变量，并返回对应的底层结构体
	ptrSubCmd := reflect.New(typeOfSubCmd)
	ptrStruct := ptrSubCmd.Elem()
	for i := 0; i < argsLen; i += 2 {
		fieldName := Str2Field[fmt.Sprintf("%s-%s", modlue, args[i])]
		// 判断对应成员是否存在
		structField, exist := typeOfSubCmd.FieldByName(fieldName)
		if !exist {
			return nil, fmt.Errorf("Cannot find field %s in %s", fieldName, typeOfSubCmd.Name())
		}
		// 设置对应的成员
		switch structField.Type.Kind() {
		case reflect.String:
			ptrStruct.FieldByName(fieldName).SetString(args[i+1])
		case reflect.Int:
			num, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Field %s need a integer value but get %s: %v", fieldName, args[i+1], err)
			}
			ptrStruct.FieldByName(fieldName).SetInt(num)
		case reflect.Interface:
			value, err := makeValue(structField.Type.Name(), args[i+1])
			if err != nil {
				return nil, fmt.Errorf("Field %s make value failed %s: %v", fieldName, args[i+1], err)
			}
			ptrStruct.FieldByName(fieldName).Set(value)
		}
	}

	return ptrSubCmd.Interface(), nil
}

func makeValue(typeName, value string) (reflect.Value, error) {
	switch typeName {
	case "Object":
		obj, err := convertToObject(value)
		return reflect.ValueOf(obj), err
	}
	return reflect.Value{}, fmt.Errorf("Unknown type: %s", typeName)
}
