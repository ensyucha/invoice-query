package dbop

import (
	"invoice/model"
)

func DeleteUncheckData(queryItem *model.Query, username string) error {

	tableName := "result_" + username

	stmt, err := db.Prepare("DELETE FROM " + tableName + " WHERE fpdm=? AND fphm=? " +
		"AND kprq=? AND totalJe=-1 AND totalSe=-1;")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(queryItem.Fpdm, queryItem.Fphm, queryItem.Kprq)

	if err != nil {
		return err
	}

	return nil
}

func UpdateFailedCheckData(queryItem *model.Query, failedState string, username string) error {

	tableName := "result_" + username

	stmt, err := db.Prepare("UPDATE " + tableName + " SET respMsg=? WHERE fpdm=? AND fphm=? " +
		"AND kprq=? AND totalJe=-1 AND totalSe=-1;")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(failedState, queryItem.Fpdm, queryItem.Fphm, queryItem.Kprq)

	if err != nil {
		return err
	}

	return nil
}