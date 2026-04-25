package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	res, err := db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", parcel.Client),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address),
		sql.Named("created_at", parcel.CreatedAt))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	parcel.Number = int(id)
	pack,err := store.Get(int(id))
	require.NoError(t, err)
	require.Equal(t, pack, parcel)
	


	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	_, err = db.Exec("delete from parcel where number = :number", sql.Named("number", id))
	require.NoError(t, err)
	_, err = store.Get(int(id))
	require.Error(t, err, "посылка не удалена")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	res, err := db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", parcel.Client),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address),
		sql.Named("created_at", parcel.CreatedAt))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)


	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	_,err = db.Exec("update parcel set address = :address where number = :number", sql.Named("address", newAddress), sql.Named("number", id))
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	pack,err := store.Get(int(id))
	require.NoError(t, err)
	require.Equal(t, newAddress, pack.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	res, err := db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", parcel.Client),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address),
		sql.Named("created_at", parcel.CreatedAt))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)


	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusSent
	_,err = db.Exec("update parcel set status = :status where number = :number", sql.Named("status", newStatus), sql.Named("number", id))
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	pack,err := store.Get(int(id))
	require.NoError(t, err)
	require.Equal(t, newStatus, pack.Status)

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	for i := 0; i < len(parcels); i++ {
		res, err := db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
			sql.Named("client", parcels[i].Client),
			sql.Named("status", parcels[i].Status),
			sql.Named("address", parcels[i].Address),
			sql.Named("created_at", parcels[i].CreatedAt))
		if err != nil {
			t.Fatal(err)
		}
		id, err := res.LastInsertId()
		require.NoError(t, err)
		require.NotZero(t, id)
		lastId := int(id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = lastId

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[lastId] = parcels[i]
	}

	// get by client // получите список посылок по идентификатору клиента, сохранённого в переменной client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		expected, ok := parcelMap[parcel.Number]
		require.True(t, ok, "посылка № %d не найдена среди добавленных", parcel.Number)
		require.Equal(t, expected, parcel)
	}
}
