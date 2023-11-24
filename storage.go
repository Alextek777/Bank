package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	MakeTransaction(sourceID int, destinationID int, amount float64) string
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createTable()
}

func (s *PostgresStore) createTable() error {
	query := `create table if not exists account (
		id 					serial primary key,
		first_name 			varchar(100),
		last_name 			varchar(100),
		number 				serial,
		encrypted_password 	varchar(100),
		balance 			serial CHECK (balance >= 0),
		created_at 			timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name, last_name, number, encrypted_password, balance, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.EncryptedPassword, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	query := fmt.Sprintf(`SELECT * FROM account WHERE number = %d`, number)
	row, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, fmt.Errorf("no Account was fount with number: %d", number)
	}

	account, err := scanIntoAccount(row)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)

	return err
}

func (s *PostgresStore) MakeTransaction(sourceID int, destinationID int, amount float64) string {

	tx, err := s.db.Begin()
	if err != nil {
		return err.Error()
	}

	query := fmt.Sprintf("UPDATE account SET balance = balance - %f	WHERE id = %d;", amount, sourceID)

	_, err = tx.Exec(query)

	if err != nil {
		return err.Error()
	}

	query = fmt.Sprintf("UPDATE account SET balance = balance + %f	WHERE id = %d;", amount, destinationID)

	_, err = tx.Exec(query)

	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return "rollback error!"
		}

		return err.Error()
	}

	err = tx.Commit()

	if err != nil {
		return err.Error()
	}

	return "Transaction complited sucessfully!"
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	querry := fmt.Sprintf(`SELECT * FROM account WHERE ID = %d`, id)
	row, err := s.db.Query(querry)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, fmt.Errorf("no Account was fount with ID: %d", id)
	}

	account, err := scanIntoAccount(row)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
