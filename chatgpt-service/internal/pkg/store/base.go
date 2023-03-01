package store

import (
	"chatgpt-service/internal/pkg/process"
	"chatgpt-service/pkg/client"
	"github.com/go-pg/pg"
)

type Database struct {
	*pg.DB
}

func (db *Database) Error() string {
	//TODO implement me
	panic("implement me")
}

func (db *Database) Heartbeat() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) StoreGptPythonSqlResult(sentence string, sql string) (int, error) {
	lastInsertId := 0
	query := "INSERT INTO flipside_query_result (query, sentence, address) VALUES (?, ?, ?) RETURNING id"
	_, err := db.QueryOne(pg.Scan(&lastInsertId), query, sql, sentence, "0x0000000000000000000000000000000000000000")
	if err != nil {
		return -1, err
	}

	return lastInsertId, nil
}

func (db *Database) UpdateCreateFlipsideQueryResult(id int, token string) error {
	// TODO: change to enum
	status := "PENDING"
	query := "UPDATE flipside_query_result SET status = ?, token = ? WHERE id = ?"
	_, err := db.Exec(query, status, token, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) StoreGetFlipsideQueryResult(request client.GetFlipsideQueryResultRequest, response client.GetFlipsideQueryResultSuccessResponse) error {
	status := "SUCCEEDED"
	query := "UPDATE flipside_query_result SET results = ?, status = ? WHERE token = ?"

	token := request.Token

	processedResults, err := process.TwoDimensionalInterfacesToJSONByte(response.Results)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, processedResults, status, token)
	if err != nil {
		return err
	}

	return nil
}
