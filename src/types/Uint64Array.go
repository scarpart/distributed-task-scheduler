package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Uint64Array []uint64

func (u *Uint64Array) Scan(value interface{}) error {
	if value == nil {
		*u = []uint64{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			*u = []uint64{}
			return nil
		}
		var arr []uint64
		if err := json.Unmarshal(v, &arr); err != nil {
			return err
		}
		*u = Uint64Array(arr)
		return nil
	case string:
		if v == "" {
			*u = []uint64{}
			return nil
		}
		var arr []uint64
		if err := json.Unmarshal([]byte(v), &arr); err != nil {
			return err
		}
		*u = Uint64Array(arr)
		return nil
	default:
		return fmt.Errorf("Cannot Scan() unsupported value type %T", value)
	}
}

func (u Uint64Array) Value() (driver.Value, error) {
	if len(u) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
