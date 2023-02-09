package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/rodaine/table"
)

func saveLogin(username, password string) error {
	os.Remove(envFile)
	return godotenv.Write(map[string]string{
		"username": username,
		"password": password,
	}, envFile)
}

func hasCsvFlag(args []string) bool {
	for i := range args {
		if args[i] == "-csv" {
			return true
		}
	}
	return false
}

func getInstaInstance() (*goinsta.Instagram, error) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("you have to login first with 'instastuff login USERNAME PASSWORD'")
	}
	instance := goinsta.New(username, password)
	if instance == nil {
		return nil, goinsta.ErrNoMore
	}
	if err := instance.Login(); err != nil {
		return nil, err
	}
	return instance, nil
}

func printTable[T any](rows []T, header []any, getRow func(row T) []any) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(header...)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i := range rows {
		tbl.AddRow(getRow(rows[i])...)
	}

	tbl.Print()
}

func saveTable[T any](rows []T, header []any, getRow func(row T) []any, name string) error {
	data := [][]string{}
	h := []string{}
	for i := range header {
		h = append(h, fmt.Sprint(header[i]))
	}
	data = append(data, h)
	for i := range rows {
		r := []string{}
		for _, v := range getRow(rows[i]) {
			r = append(r, fmt.Sprint(v))
		}
		data = append(data, r)
	}

	csvFile, err := os.Create(fmt.Sprintf("%v-%v.csv", name, strings.ReplaceAll(time.Now().Local().String(), " ", "_")))

	if err != nil {
		return err
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	w.Comma = ';'
	w.WriteAll(data)
	w.Flush()
	return w.Error()
}
