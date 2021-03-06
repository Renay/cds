package build

import (
	"database/sql"
	"encoding/json"

	"github.com/ovh/cds/engine/api/database"
	"github.com/ovh/cds/sdk"
)

// LoadTestResults retrieves tests on a specific build in database
func LoadTestResults(db *sql.DB, pbID int64) (sdk.Tests, error) {
	query := `SELECT tests FROM pipeline_build_test WHERE pipeline_build_id = $1`
	t := sdk.Tests{}
	var data string

	err := db.QueryRow(query, pbID).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, nil
		}
		return t, err
	}

	err = json.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// InsertTestResults inserts test results of a specific pipeline build in database
func InsertTestResults(db database.Executer, pbID int64, tests sdk.Tests) error {
	query := `INSERT INTO pipeline_build_test (pipeline_build_id, tests) VALUES ($1, $2)`

	data, err := json.Marshal(tests)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, pbID, string(data))
	if err != nil {
		return err
	}

	return nil
}

// UpdateTestResults update test results of a specific pipeline build in database
func UpdateTestResults(db *sql.DB, pbID int64, tests sdk.Tests) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM pipeline_build_test WHERE pipeline_build_id = $1`
	_, err = tx.Exec(query, pbID)
	if err != nil {
		return err
	}

	err = InsertTestResults(tx, pbID, tests)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// DeleteTestResults removes from database test results for a specific pipeline build
func DeleteTestResults(db database.Executer, pbID int64) error {
	query := `DELETE FROM pipeline_build_test WHERE pipeline_build_id = $1`

	_, err := db.Exec(query, pbID)
	if err != nil {
		return err
	}

	return nil
}

// DeletePipelineTestResults removes from database test results for a specific pipeline
func DeletePipelineTestResults(db database.Executer, pipID int64) error {
	query := `DELETE FROM pipeline_build_test WHERE pipeline_build_id IN
		(SELECT id FROM pipeline_build WHERE pipeline_id = $1)`

	_, err := db.Exec(query, pipID)
	if err != nil {
		return err
	}

	return nil
}

/*
// DeleteApplicationPipelineTestResults removes from database test results for a specific pipeline linked to a specific application
func DeleteApplicationPipelineTestResults(db database.Executer, appID int64, pipID int64) error {
	query := `DELETE FROM pipeline_build_test WHERE pipeline_build_id IN
		(SELECT id FROM pipeline_build WHERE application_id = $1 AND pipeline_id = $2)`

	_, err := db.Exec(query, appID, pipID)
	if err != nil {
		return err
	}

	return nil
}
*/
