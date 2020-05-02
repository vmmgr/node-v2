package db

import (
	"fmt"
	"log"
)

//VM

func AddDBVM(data VM) bool {
	fmt.Println("add database: " + data.Name)
	db := *connectdb()
	defer db.Close()

	addDb, err := db.Prepare(`INSERT INTO "vm" ("name","cpu","memory","storage","storagepath","net","vnc","socket","status","autostart") VALUES (?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: SQL Prepare Error")
		return false
	}

	if _, err := addDb.Exec(data.Name, data.CPU, data.Mem, data.Storage, data.StoragePath, data.Net, data.Vnc, data.Socket, data.Status, data.AutoStart); err != nil {
		fmt.Println(err)
		fmt.Println("Error: SQL Exec Error")
		return false
	}
	return true
}

func DeleteDBVM(id int) bool {
	db := connectdb()
	defer db.Close()

	deleteDb := "DELETE FROM vm WHERE id = ?"
	_, err := db.Exec(deleteDb, id)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: SQL open Error")
		return false
	}
	return true
}

func VMDBGetAll() []VM {
	db := *connectdb()
	defer db.Close()

	cmd := "SELECT * FROM vm"
	rows, _ := db.Query(cmd)

	defer rows.Close()

	var bg []VM
	for rows.Next() {
		var b VM
		err := rows.Scan(&b.ID, &b.Name, &b.CPU, &b.Mem, &b.Storage, &b.StoragePath, &b.Net, &b.Vnc, &b.Socket, &b.Status, &b.AutoStart)
		if err != nil {
			fmt.Println(err)
		}
		bg = append(bg, b)
	}
	return bg
}

func VMDBGetVMID(name string) (int, error) {
	data := VMDBGetAll()
	for i, _ := range data {
		if data[i].Name == name {
			return data[i].ID, nil
		}
	}

	return -1, fmt.Errorf("Not Found!!!")
}

func VMDBGetVMStatus(id int) (int, error) {
	//0: PowerOff 1: PowerOn 2:Suspend 3: TmpStop 4: busy
	data := VMDBGetAll()
	for i, _ := range data {
		if data[i].ID == id {
			return data[i].Status, nil
		}
	}

	return -1, fmt.Errorf("Not Found!!!")
}

func VMDBStatusUpdate(id, status int) bool {
	db := *connectdb()
	defer db.Close()

	cmd := "UPDATE vm SET status = ? WHERE id = ?"
	_, err := db.Exec(cmd, status, id)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func VMDBGetData(id int) (*VM, error) {
	data := VMDBGetAll()
	var result VM
	for i, _ := range data {
		if data[i].ID == id {
			result = data[i]
			fmt.Println(i)
			return &result, nil
		}
	}
	return &result, fmt.Errorf("Not Found")
}

func VMDBUpdate(data *VM) {
}

func VMDBStatusStop(id int) {

}
