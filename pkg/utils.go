package pkg

import (
	"database/sql"
	"math"
	"time"
)

func StringToNullString(value string) sql.NullString {
	if len(value) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func NullStringToString(value sql.NullString) string {
	if !value.Valid {
		return "null"
	}

	return value.String
}

func TimeToNullTime(value time.Time) sql.NullTime {
	nullTime := time.Time{}

	if value == nullTime {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  value,
		Valid: true,
	}
}

func NullTimeToString(value sql.NullTime) string {
	if !value.Valid {
		return "null"
	}

	return value.Time.Format(time.RFC3339)
}

func GetDroneDistance(x float64, y float64) int64 {
	const XC float64 = 250000
	const YC float64 = 250000

	return int64(math.Round(math.Sqrt(math.Pow(x-XC, 2) + math.Pow(y-YC, 2))))
}
