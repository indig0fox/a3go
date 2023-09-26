package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/indig0fox/a3go/a3interface"
)

func receiveTestCommand(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {

	dataSlice := strings.Split(data, "|")
	return fmt.Sprintf(`["Called by %s", "%s", "%s", "%s"]`,
		ctx.SteamID,
		dataSlice[0],
		dataSlice[1],
		dataSlice[2],
	), nil
}

func receiveTestCommandArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {

	return fmt.Sprintf(`["Called by %s", "%s", "%s", "%s"]`,
		ctx.SteamID,
		command,
		args[0],
		args[1],
	), nil
}

func saveCaller(ctx a3interface.ArmaExtensionContext, data string) (string, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Create table if it doesn't exist
	sqlStmt := `
	create table if not exists call_log (id integer not null primary key, player_uid text, server_name text, mission_name text, file_source text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return "", err
	}

	// Insert data
	stmt, err := db.Prepare("INSERT INTO call_log(player_uid, server_name, mission_name, file_source) values(?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ctx.SteamID, ctx.ServerName, ctx.MissionNameSource, ctx.FileSource)
	if err != nil {
		return "", err
	}

	return "Logged row to database!", nil
}
