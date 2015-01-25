package main

import (
	`net/http`
	mgo `gopkg.in/mgo.v2`
	`fmt`
)

func searchHandler(dbConn *mgo.Session, dbName string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(writer, `%s`, dbName)
	}
}
