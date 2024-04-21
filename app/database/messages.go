package database

func GetMessageCount(unread bool) (uint16, error) {
	var count uint16
	var statement string
	if unread {
		statement = `select count(*) from messages where read = FALSE;`
	} else {
		statement = `select count(*) from messages;`
	}
	rows, err := db.Query(statement)
	if err != nil {
		return count, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func GetAllMessages() ([]Message, error) {
	var messages []Message
	statement := `select type, header, message, sent, read, id from messages;`

	rows, err := db.Query(statement)
	if err != nil {
		return []Message{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.Type, &msg.Header, &msg.Message, &msg.Sent, &msg.Read, &msg.Id)
		if err != nil {
			return []Message{}, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func MessageRead(id uint64) error {
	statement := `update messages SET read = $2 where id = $1`
	_, err := db.Exec(statement, id, true)
	return err
}

func DeleteMessage(id uint64) error {
	statement := `delete from messages where id = $1`
	_, err := db.Exec(statement, id)
	return err
}
