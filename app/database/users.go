package database

func CreateUser(user *User) error {
	statement := `insert into users(name, email, password) values($1, $2, $3);`
	// db.Create(&user)
	_, err := db.Exec(statement, user.Name, user.Email, user.Password)
	return err
}

func CheckEmail(email string, user *User) bool {
	statement := `select id, name, email, password from users where email=$1 limit 1;`
	rows, err := db.Query(statement, email)
	if err != nil {
		return false
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return false
		}
	}
	return true
}

func GetUser(id string) (User, error) {
	var user User
	statement := `select * from users where id=$1;`
	rows, err := db.Query(statement, id)
	if err != nil {
		return User{}, err
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}

func GetUsers() ([]User, error) {
	var users []User

	statement := `select id, name, email, password from users;`

	rows, err := db.Query(statement)
	if err != nil {
		return []User{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var name, email, password string
		var id uint64

		err = rows.Scan(&id, &name, &email, &password)
		if err != nil {
			return []User{}, err
		}

		user := User{
			ID:       id,
			Name:     name,
			Email:    email,
			Password: password,
		}

		users = append(users, user)
	}

	return users, nil
}

func DeleteUser(id string) error {
	statement := `delete users where id=$1;`
	_, err := db.Exec(statement, id)
	return err
}

func ResetPassword(id uint64, password string) error {
	statement := `update users set password=$2 where id=$1;`
	_, err := db.Exec(statement, id, password)
	return err
}
