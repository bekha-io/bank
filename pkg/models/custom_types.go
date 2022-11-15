package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type RawJson json.RawMessage

func (r *RawJson) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type. Cannot insert into the database")
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*r = RawJson(result)
	return err
}

func (r RawJson) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.RawMessage(r).MarshalJSON()
}
