package db

//Add
func AddDBStorage(data Storage) result {
	db := InitDB()
	defer db.Close()
	//db.Table("group").CreateTable(&data)
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Result: false, ID: 0}
	} else {
		return result{Result: true, ID: data.ID}
	}
}

//Delete
func DeleteDBStorage(data Storage) bool {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Update
func UpdateDBStorage(data Storage) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Storage{GroupID: data.GroupID, Name: data.Name, Driver: data.Driver, MaxSize: data.MaxSize, Lock: data.Lock})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBStorage() []Storage {
	db := InitDB()
	defer db.Close()

	var vm []Storage
	db.Find(&vm)
	return vm
}

func SearchDBStorage(data Storage) Storage {
	db := InitDB()
	defer db.Close()

	var result Storage
	//search StorageName and StorageID
	if data.Name != "" {
		db.Where("name = ?", data.Name).First(&result)
	} else if data.ID != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("ID = ?", data.ID).First(&result)
	}

	return result
}
