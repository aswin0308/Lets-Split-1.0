package mysql

import (
	"expense/pkg/models"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string,) error {

	log.Println(m.DB);
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

	stmt := `INSERT INTO user (name, email,password)
				VALUES(?,?,?)`

	_,err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil

}

func (m *UserModel) Authenticate(email, password string) (bool, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return false, models.ErrInvalidCredentials
	} else if err != nil {
		return true, err
	}

	//compare the provided password with the hashped password. If they match
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err == bcrypt.ErrMismatchedHashAndPassword {
        return false, models.ErrInvalidCredentials
    } else if err != nil {
        return true, err
    }

	return false, nil
}

func (m *UserModel) GetAllUsers() ([]*models.User, error) {
	stmt := `SELECT userId,name from user`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	UsersinDB := []*models.User{}

	for rows.Next() {
		// Create a pointer to a new zeroed Todos struct.
		s := &models.User{}

		err = rows.Scan(&s.UserID, &s.Name)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of todos.
		UsersinDB = append(UsersinDB, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return UsersinDB, nil
}