package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/cristalhq/acmd"
	"github.com/joho/godotenv"
)

const envFile = "instastuff.env"

var noEnv = true

func init() {
	err := godotenv.Load(envFile)
	if err != nil {
		return
	}
	noEnv = false
}

func main() {
	// argsWithoutProg := os.Args[1:]

	// insta := goinsta.New("USERNAME", "PASSWORD")

	// // Export your configuration
	// // after exporting you can use Import function instead of New function.
	// // insta, err := goinsta.Import("~/.goinsta")
	// // it's useful when you want use goinsta repeatedly.
	// insta.Export("~/.goinsta")

	cmds := []acmd.Command{
		{
			Name:        "login",
			Description: "saves the login data for future actions",
			ExecFunc:    cmdLogin,
		},
		{
			Name:        "followers",
			Description: "Shows all followers of the loggedin account (use -csv to download it as csv file, eg: instastuff followers -csv)",
			ExecFunc:    cmdFollowers,
		},
		{
			Name:        "following",
			Description: "Shows all people that are followed by the loggedin account (use -csv to download it as csv file, eg: instastuff following -csv)",
			ExecFunc:    cmdFollowing,
		},
	}

	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:         "instastuff",
		AppDescription:  "Helpers for insta",
		PostDescription: "Support gibs ned",
		Version:         "0.0.1",
		Output:          os.Stdout,
		Args:            os.Args,
	})

	if err := r.Run(); err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("\n")
}

func cmdLogin(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("you must provide a username and a password")
	}
	username := args[0]
	password := args[1]
	return saveLogin(username, password)
}

func cmdFollowers(ctx context.Context, args []string) error {
	insta, err := getInstaInstance()
	if err != nil {
		return err
	}
	fi := insta.Account.Followers()
	followers := fi.Users
	for fi.Next() {
		followers = append(followers, fi.Users...)
	}
	if hasCsvFlag(args) {
		err := saveTable(followers, []any{"ID", "Name"}, func(row goinsta.User) []any {
			return []any{row.ID, row.FullName}
		}, "followers")
		fmt.Printf("CSV saved successfully")
		return err
	} else {
		printTable(followers, []any{"ID", "Name"}, func(row goinsta.User) []any {
			return []any{row.ID, row.FullName}
		})
		fmt.Printf("\nYou have %v followers", len(followers))
	}
	return nil
}

func cmdFollowing(ctx context.Context, args []string) error {
	insta, err := getInstaInstance()
	if err != nil {
		return err
	}
	fi := insta.Account.Following()
	followers := fi.Users
	for fi.Next() {
		followers = append(followers, fi.Users...)
	}
	if hasCsvFlag(args) {
		err := saveTable(followers, []any{"ID", "Name"}, func(row goinsta.User) []any {
			return []any{row.ID, row.FullName}
		}, "following")
		fmt.Printf("CSV saved successfully")
		return err
	} else {
		printTable(followers, []any{"ID", "Name"}, func(row goinsta.User) []any {
			return []any{row.ID, row.FullName}
		})
		fmt.Printf("\nYou follow %v people", len(followers))
	}
	return nil
}
