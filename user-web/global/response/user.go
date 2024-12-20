package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	temp := time.Time(j).Format("2006-01-02 15:04:05")
	stmp := fmt.Sprintf("\"%s\"", temp)

	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Nickname string `json:"name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Mobile   string   `json:"mobile"`
	Gender   string   `json:"gender"`
}
