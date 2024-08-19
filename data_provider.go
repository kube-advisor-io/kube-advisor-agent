package main 

type DataProvider interface {
	GetData() map[string]interface{}
}