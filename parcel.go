package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	answer, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt) // добавляем новую строку в таблицу
	if err != nil {
		return 0, err // возвращаем ошибку
	}
	lastID, err := answer.LastInsertId() // идентификатор последней добавленной записи
	if err != nil {
		return 0, err
	}
	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	answer := s.db.QueryRow("SELECT * FROM parcel WHERE number = ?", number) // запрос одной строки
	err := answer.Err()                                                      // возвращаем ошибку
	if err != nil {
		return Parcel{}, err // возвращаем ошибку
	}
	p := Parcel{}
	err = answer.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	answer, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client) // запрос множества строк
	if err != nil {
		return nil, err // возвращаем ошибку
	}
	var res []Parcel // создадим срез Parcel
	for answer.Next() {
		p := Parcel{}
		err = answer.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt) // заполним объект данными
		if err != nil {
			return nil, err // возвращаем ошибку
		}
		res = append(res, p)
	}
	if err = answer.Err(); err != nil {
		return nil, err // возвращаем ошибку
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number) // обновляем статус
	if err != nil {
		return err // возвращаем ошибку
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number) // получим строку по идентификатору
	if err != nil {
		return err // возвращаем ошибку
	}
	if parcel.Status != ParcelStatusRegistered {
		return errors.New("Cannot change address. Parcel status is not registered.") // возвращаем ошибку
	}
	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number) // обновляем адрес
	if err != nil {
		return err // возвращаем ошибку
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	parsel, err := s.Get(number) // получим строку по идентификатору
	if err != nil {
		return err // возвращаем ошибку
	}
	if parsel.Status != ParcelStatusRegistered {
		return errors.New("Cannot delete parcel. Parcel status is not registered.") // возвращаем ошибку
	}
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number) // удалим строку
	if err != nil {
		return err // возвращаем ошибку
	}
	return nil
}
