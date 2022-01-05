package types

import (
	"encoding/json"
	"errors"
	"time"
)

//Duration 是对time.Duration的封装，主要解决json序列化的问题
type Duration struct {
	time.Duration
}

//MarshalJSON 实现json接口
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

//UnmarshalJSON 实现json接口
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
