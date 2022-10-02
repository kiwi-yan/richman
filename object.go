package main

type ObjType = int

const (
	objFood ObjType = iota
)

type Object interface {
	objType() ObjType
}

type Food struct{}

func (Food) objType() ObjType { return objFood }

var PriceList = []struct {
	ObjType
	price int
}{
	{objFood, 10},
}
