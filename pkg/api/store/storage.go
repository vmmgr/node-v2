package store

//Add
func AddDBStorage(data Storage) result {
	db := InitDB()
	defer db.Close()
	//store.Table("group").CreateTable(&gateway)
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil, ID: data.ID}
}

//Delete
func DeleteDBStorage(data Storage) result {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	} else {
		return result{Error: nil, ID: data.ID}
	}
}

//Update
func UpdateDBStorage(data Storage) result {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Storage{GroupID: data.GroupID, Name: data.Name, Driver: data.Driver, MaxSize: data.MaxSize, Lock: data.Lock})

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil, ID: data.ID}
}

//Get
func SearchDBStorage(data Storage) (Storage, error) {
	db := InitDB()
	defer db.Close()
	var result Storage

	db.Where("ID = ?", data.ID).First(&result)
	if err := db.Error; err != nil {
		return Storage{}, err
	}
	return result, nil
}

func GetAllDBStorage() ([]Storage, error) {
	db := InitDB()
	defer db.Close()

	var vm []Storage
	db.Find(&vm)

	if err := db.Error; err != nil {
		return []Storage{}, err
	}
	return vm, nil
}
