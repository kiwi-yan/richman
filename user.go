package main

import "fmt"

const InitMoney = 1000

type User struct {
	name    string
	account int
	bag     map[Object]int
	assets  []Assets
	energy  int
}

func NewUser(name string) *User {
	return &User{
		name:    name,
		account: InitMoney,
		bag:     make(map[Object]int),
		assets:  make([]Assets, 0),
		energy:  100,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("User %s's account has %d coins and %#v", u.name, u.account, u.bag)
}

func (u *User) Buy(obj Object, num int) bool {
	if PriceList[obj.objType()].price*num <= u.account {
		u.account -= PriceList[obj.objType()].price * num
		u.bag[obj] += num
		return true
	}
	return false
}

func (u *User) Sell(obj Object, num int) bool {
	if u.bag[obj] >= num {
		u.bag[obj] -= num
		u.account += PriceList[obj.objType()].price * num
		if u.bag[obj] == 0 {
			delete(u.bag, obj)
		}
		return true
	}
	return false
}
