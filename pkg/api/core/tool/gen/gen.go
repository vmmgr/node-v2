package gen

import (
	"fmt"
	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	uu := u.String()

	return uu, nil
}
