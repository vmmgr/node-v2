package vm

import "github.com/vmmgr/node/db"

func vmExistsName(name string) bool {
	if data, err := db.GetAllDBVM(); err != nil {
		return false
	} else {
		for _, v := range data {
			if v.Name == name {
				return false
			}
		}
		return true
	}
}

func vmExistsID(id int) bool {
	if data, err := db.GetAllDBVM(); err != nil {
		return false
	} else {
		for _, v := range data {
			if v.ID == id {
				return false
			}
		}
		return true
	}
}
