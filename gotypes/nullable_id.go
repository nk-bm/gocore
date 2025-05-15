package gotypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NullIdentityID represents an IdentityID that may be null.
// NullIdentityID implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullInt64.
type NullIdentityID struct {
	IdentityID IdentityID
	Valid      bool // Valid is true if IdentityID is not NULL
}

// Scan implements the Scanner interface for NullIdentityID
func (ni *NullIdentityID) Scan(value interface{}) error {
	if value == nil {
		ni.IdentityID, ni.Valid = 0, false
		return nil
	}

	switch v := value.(type) {
	case int64:
		ni.IdentityID = IdentityID(v)
		ni.Valid = true
	case []byte:
		var i int64
		if err := json.Unmarshal(v, &i); err != nil {
			return err
		}
		ni.IdentityID = IdentityID(i)
		ni.Valid = true
	default:
		return fmt.Errorf("cannot convert %T to NullIdentityID", value)
	}
	return nil
}

// Value implements the driver Valuer interface for NullIdentityID
func (ni NullIdentityID) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return int64(ni.IdentityID), nil
}

// MarshalJSON implements the json.Marshaler interface for NullIdentityID
func (ni NullIdentityID) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.IdentityID)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullIdentityID
func (ni *NullIdentityID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}
	var i IdentityID
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	ni.IdentityID = i
	ni.Valid = true
	return nil
}
