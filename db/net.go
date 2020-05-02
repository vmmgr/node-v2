package db

//Add
func AddDBNet(data Net) result {
	db := InitDB()
	defer db.Close()
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Result: false, ID: 0}
	} else {
		return result{Result: true, ID: data.ID}
	}
}

//Delete
func DeleteDBNet(data Net) bool {
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
func UpdateDBNet(data Net) bool {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Net{Name: data.Name, Driver: data.Driver, Vlan: data.Vlan, Status: data.Status})

	if err := db.Error; err != nil {
		db.Rollback()
		return false
	} else {
		return true
	}
}

//Get
func GetAllDBNet() []Net {
	db := InitDB()
	defer db.Close()

	var vm []Net
	db.Find(&vm)
	return vm
}

func SearchDBNet(data Net) Net {
	db := InitDB()
	defer db.Close()

	var result Net
	//search NetName and NetID
	if data.Name != "" {
		db.Where("name = ?", data.Name).First(&result)
	} else if data.ID != 0 { //初期値0であることが前提　確認の必要あり
		db.Where("ID = ?", data.ID).First(&result)
	}

	return result
}
