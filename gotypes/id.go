package gotypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IdentityID int64

func (id IdentityID) IsPositive() bool {
	return id < 0
}

func (id IdentityID) IsNegative() bool {
	return id > 0
}

func IntToIdentityID(id int, isArtist bool) IdentityID {
	if isArtist {
		return IdentityID(-int64(id))
	}
	return IdentityID(id)
}

func Int64ToIdentityID(id int64, isArtist bool) IdentityID {
	if isArtist {
		return IdentityID(-id)
	}
	return IdentityID(id)
}

func (id IdentityID) ToInt64() int64 {
	if id < 0 {
		return -int64(id)
	}
	return int64(id)
}

// Scan implements the Scanner interface for IdentityID
func (id *IdentityID) Scan(value any) error {
	if value == nil {
		*id = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*id = IdentityID(v)
	case []byte:
		var i int64
		if err := json.Unmarshal(v, &i); err != nil {
			return err
		}
		*id = IdentityID(i)
	default:
		return fmt.Errorf("cannot convert %T to IdentityID", value)
	}
	return nil
}

// Value implements the driver Valuer interface for IdentityID
func (id IdentityID) Value() (driver.Value, error) {
	return int64(id), nil
}

// MarshalJSON implements the json.Marshaler interface for IdentityID
func (id IdentityID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(id))
}

// UnmarshalJSON implements the json.Unmarshaler interface for IdentityID
func (id *IdentityID) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*id = IdentityID(i)
	return nil
}
