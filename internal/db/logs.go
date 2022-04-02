package db

import (
	"context"
	"fmt"
	"github.com/ashirko/imlogs/internal/models"
	"github.com/codenotary/immudb/pkg/errors"
	"log"
	"reflect"
)

func (db Database) GetLogs(opts map[string]interface{}) ([]models.LogLine, error) {
	logLines := []models.LogLine{}

	// Create initial query, then add optional conditions and limitations.
	query := "SELECT id, create_time, source, line FROM logs "

	if _, ok := opts["source"]; ok {
		query += "WHERE source = @source "
	}

	if count, ok := opts["count"]; ok {
		// TODO the line below doesn't work for some reason. Parameterized LIMIT is not supported?
		// query += "ORDER BY id DESC LIMIT @count  "
		if reflect.TypeOf(count).Kind() == reflect.Int {
			query = fmt.Sprintf(query+"ORDER BY id DESC LIMIT %d", count)
		} else {
			return logLines, errors.New("incorrect value of count")
		}
	}

	// Execute the query.
	result, err := db.Client.SQLQuery(context.TODO(), query, opts, true)
	if err != nil {
		return logLines, err
	}

	// Build the result.
	for _, row := range result.Rows {
		logLine := models.NewLogLine(row.Values[2].GetS(), row.Values[3].GetS())
		logLines = append(logLines, *logLine)
	}

	return logLines, nil
}

func (db Database) AddLogLine(logLine models.LogLine) error {
	_, err := db.Client.SQLExec(context.TODO(), `
		INSERT INTO logs (create_time, source, line)
		VALUES (NOW(), @source, @line)`,
		map[string]interface{}{"source": logLine.Source, "line": logLine.Line})
	return err
}

func (db Database) AddLogLines(logLines models.LogLines) error {
	// Create a new transaction.
	tx, err := db.Client.NewTx(context.TODO())
	if err != nil {
		return err
	}

	// Add each new log line in a separated query.
	for _, l := range logLines.Logs {
		err = tx.SQLExec(context.TODO(), `
			INSERT INTO logs (create_time, source, line)
			VALUES (NOW(), @source, @line)`,
			map[string]interface{}{"source": l.Source, "line": l.Line})
		if err != nil {
			if err1 := tx.Rollback(context.TODO()); err1 != nil {
				log.Println("Rollback failed: ", err1)
			}
			return err
		}
	}
	_, err = tx.Commit(context.TODO())
	return err
}

func (db Database) GetLogsCount() (int64, error) {
	result, err := db.Client.SQLQuery(context.TODO(), "SELECT COUNT(*) FROM logs", map[string]interface{}{}, true)
	if err != nil {
		return 0, err
	}
	count := result.Rows[0].Values[0].GetN()
	return count, nil
}

func (db Database) createLogsTable() error {
	_, err := db.Client.SQLExec(context.TODO(), `
		CREATE TABLE IF NOT EXISTS logs (
			id           INTEGER AUTO_INCREMENT,
			create_time  TIMESTAMP,
			source       VARCHAR[60],
			line         VARCHAR NOT NULL,
			PRIMARY KEY id)
		`, map[string]interface{}{})
	return err
}
