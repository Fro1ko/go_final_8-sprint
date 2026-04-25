package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	lastID := int(id)
	return lastID, nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	p := Parcel{}
	err := s.db.QueryRow("select number, client, status, address, created_at from parcel where number = :number", sql.Named("number", number)).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	// заполните объект Parcel данными из таблицы

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	rows,err := s.db.Query("select number, client, status, address, created_at from parcel where client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("update parcel set status = :status where number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	
	_, err := s.db.Exec("update parcel set address = :address where number = :number and status = :status", sql.Named("address", address), sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	
	_,err := s.db.Exec("delete from parcel where number = :number and status = :status", sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	
	return nil
}
