package structure

import (
	"encoding/json"
	"github.com/jinzhu/copier"
)

// Copy 结构体映射
func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}

// Copy2
func Copy2(s, ts interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, ts)
	return err
}
