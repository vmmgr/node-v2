package db

import "fmt"

//Add
func AddDBVM(data VM) result {
	db := InitDB()
	defer db.Close()
	//db.Table("group").CreateTable(&data)
	db.Create(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil, ID: data.ID}
}

//Delete
func DeleteDBVM(data VM) result {
	db := InitDB()
	defer db.Close()
	db.Delete(&data)

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err, ID: 0}
	}
	return result{Error: nil, ID: data.ID}
}

//Update
func UpdateDBVM(data VM) result {
	db := InitDB()
	defer db.Close()
	db.Model(&data).Updates(VM{CPU: data.CPU, Mem: data.Mem, Status: data.Status, AutoStart: data.AutoStart})

	if err := db.Error; err != nil {
		db.Rollback()
		return result{Error: err}
	}
	return result{Error: nil, ID: data.ID}
}

//Get

func SearchDBVM(data VM) (VM, error) {
	db := InitDB()
	defer db.Close()

	var result VM
	db.Where("ID = ?", data.ID).First(&result)

	if err := db.Error; err != nil {
		return data, fmt.Errorf("Error: DB Error ")
	}
	return result, nil
}

func GetAllDBVM() ([]VM, error) {
	db := InitDB()
	defer db.Close()
	var vm []VM
	db.Find(&vm)

	if err := db.Error; err != nil {
		return []VM{}, err
	}
	return vm, nil
}
