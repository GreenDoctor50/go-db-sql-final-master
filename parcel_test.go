package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	padd, err := store.Add(parcel) // добавляем новую строку в таблицу
	require.NoError(t, err)        // возвращаем ошибку
	require.NotZero(t, padd)       // идентификатор последней добавленной записи

	pget, err := store.Get(padd) // получаем только что добавленную посылку
	require.NoError(t, err)      // возвращаем ошибку
	require.Equal(t, padd, pget) // проверяем, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel

	err = store.Delete(padd) // удаляем добавленную посылку
	require.NoError(t, err)  // возвращаем ошибку

	_, err = store.Get(padd) // попытка получить удаленную посылку
	require.Error(t, err)    // возвращаем ошибку
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)                     // возвращаем ошибку
	defer db.Close()                            // закроем ресурс после использования

	store := NewParcelStore(db)
	parsel := getTestParcel()

	p, err := store.Add(parsel) // добавляем новую посылку в БД
	require.NoError(t, err)     // возвращаем ошибку
	require.NotEmpty(t, p)      // идентификатор последней добавленной записи

	newAddressForUpdate := "new test address"
	err = store.SetAddress(p, newAddressForUpdate) // Обновляем адрес
	require.NoError(t, err)                        // возвращаем ошибку

	pget, err := store.Get(p)                           // получаем только что добавленную посылку
	require.NoError(t, err)                             // возвращаем ошибку
	require.Equal(t, newAddressForUpdate, pget.Address) // проверяем адрес
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)                     // возвращаем ошибку
	defer db.Close()                            // закроем ресурс после использования

	store := NewParcelStore(db)
	parsel := getTestParcel()

	idp, err := store.Add(parsel) // добавляем новую посылку в БД
	require.NoError(t, err)       // возвращаем ошибку
	require.NotEmpty(t, idp)      // идентификатор последней добавленной записи

	err = store.SetStatus(idp, ParcelStatusSent) // Обновляем статус
	require.NoError(t, err)                      // возвращаем ошибку

	p, err := store.Get(idp)                     // получаем только что добавленную посылку
	require.NoError(t, err)                      // возвращаем ошибку
	require.Equal(t, ParcelStatusSent, p.Status) // проверяем статус
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)                     // возвращаем ошибку
	defer db.Close()                            // закроем ресурс после использования

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

	for i := 0; i < len(parcels); i++ {
		id, err := db.Add(parcels[i]) // добавляем новую посылку в БД
		require.NoError(t, err)       // возвращаем ошибку
		require.NotEmpty(t, id)       // идентификатор последней добавленной записи

		parcels[i].Number = id     // обновляем идентификатор добавленной у посылки
		parcelMap[id] = parcels[i] // сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
	}

	store := NewParcelStore(db)
	storedParcels, err := store.GetByClient(client)   // список посылок по идентификатору клиента, сохранённого в переменной client
	require.NoError(t, err)                           // возвращаем ошибку
	assert.Equal(t, len(parcels), len(storedParcels)) // убеждаемся, что количество полученных посылок совпадает с количеством добавленных

	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
