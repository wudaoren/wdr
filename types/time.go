package types

import (
	"time"
)

//自定义时间类型
type Time time.Time

// MarshalJSON satify the json marshal interface
func (t *Time) MarshalJSON() ([]byte, error) {
	if time.Time(*t).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(*t).Format("2006-01-02 15:04:05") + `"`), nil
}

//
func (t *Time) IsZero() bool {
	return time.Time(*t).IsZero()
}

//
func (t *Time) Format() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

// returns time.Now() no matter what!
func (t *Time) UnmarshalJSON(data []byte) error {
	val := string(data)
	if val == `""` {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, val, time.Local)
	if err != nil {
		return err
	}
	*t = Time(now)
	return nil
}
