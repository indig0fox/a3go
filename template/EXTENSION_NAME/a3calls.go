package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/indig0fox/a3go/a3interface"
	"github.com/indig0fox/a3go/assemblyfinder"
	_ "github.com/mattn/go-sqlite3"
)

func ReceiveTestCommand(
	ctx a3interface.ArmaExtensionContext,
	data string,
) (string, error) {

	dataSlice := strings.Split(data, "|")
	dataSliceWithoutPrefix := dataSlice[1:]
	for i, v := range dataSliceWithoutPrefix {
		dataSliceWithoutPrefix[i] = fmt.Sprintf(`%q`, v)
	}

	s := fmt.Sprintf(`["Called by %s", [%s]]`,
		ctx.SteamID,
		strings.Join(dataSliceWithoutPrefix, ", "),
	)
	fmt.Println(s)

	return s, nil
}

func ReceiveTestCommandArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {

	return fmt.Sprintf(`["Called by %s", %q, %q]`,
		ctx.SteamID,
		command,
		args,
	), nil
}

func ReturnJSONFromHashMapArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {

	JSONInterface, err := a3interface.ParseSQF(args[0])
	if err != nil {
		return "", err
	}
	JSONMapStringInterface, err := a3interface.ParseSQFHashMap(JSONInterface)
	if err != nil {
		return "", err
	}

	JSONString, err := json.MarshalIndent(JSONMapStringInterface, "", "  ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s`, JSONString), nil
}

func SaveCallerArgs(
	ctx a3interface.ArmaExtensionContext,
	command string,
	args []string,
) (string, error) {
	_, err := SaveCaller(ctx, args[0])
	if err != nil {
		return "", err
	}

	arrArg1, err := a3interface.ParseSQF(args[0])
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf(`["Logged row!", Args: %q, Parsed: %q]`, args, arrArg1)
	return res, nil

}

func SaveCaller(ctx a3interface.ArmaExtensionContext, data string) (string, error) {
	modulePath := assemblyfinder.GetModulePath()
	moduleDir := filepath.Dir(modulePath)
	fmt.Println("moduleDir: ", moduleDir)
	db, err := sql.Open("sqlite3", filepath.Join(moduleDir, "call_log.db"))
	if err != nil {
		fmt.Println(err)
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

	return `["Logged row to database!"]`, nil
}
