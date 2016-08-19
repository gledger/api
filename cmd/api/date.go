package main

import (
	"encoding/json"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", v)
	*d = Date(t)
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	return json.Marshal(t.Format("2006-01-02"))
}
