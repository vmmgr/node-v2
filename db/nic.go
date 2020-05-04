package db

//Add
func AddDBNIC(data NIC) result {
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
func DeleteDBNIC(data NIC) result {
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
func UpdateDBNIC(data NIC) result {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(NIC{Name: data.Name, Driver: data.Driver, GroupID: data.GroupID, NetID: data.NetID, MacAddress: data.MacAddress, Lock: data.Lock})

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil}
}

//Get
func SearchDBNIC(data NIC) (NIC, error) {
	db := InitDB()
	defer db.Close()
	var result NIC

	db.Where("ID = ?", data.ID).First(&result)

	if err := db.Error; err != nil {
		return NIC{}, err
	}
	return result, nil
}

func GetAllDBNIC() ([]NIC, error) {
	db := InitDB()
	defer db.Close()

	var net []NIC
	db.Find(&net)
	if err := db.Error; err != nil {
		return []NIC{}, err
	}
	return net, nil
}
