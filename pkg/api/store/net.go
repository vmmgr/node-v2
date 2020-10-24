package store

//Add
func AddDBNet(data Net) result {
	db := InitDB()
	defer db.Close()
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil, ID: data.ID}
}

//Delete
func DeleteDBNet(data Net) result {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil}
}

//Update
func UpdateDBNet(data Net) result {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(Net{Name: data.Name, GroupID: data.GroupID, VLAN: data.VLAN, Lock: data.Lock})

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil}
}

//Get
func SearchDBNet(data Net) (Net, error) {
	db := InitDB()
	defer db.Close()
	var result Net

	db.Where("ID = ?", data.ID).First(&result)

	if err := db.Error; err != nil {
		return Net{}, err
	}
	return result, nil
}

func GetAllDBNet() ([]Net, error) {
	db := InitDB()
	defer db.Close()

	var net []Net
	db.Find(&net)
	if err := db.Error; err != nil {
		return []Net{}, err
	}
	return net, nil
}
