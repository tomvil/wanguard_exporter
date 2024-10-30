package collectors

import (
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
)

func bitsToBytes(b float64) float64 {
	return b / 8
}

func getFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case nil:
		return 0
	case string:
		r := strings.NewReplacer(
			" days", "",
			"∞", "9999")

		result, err := strconv.ParseFloat(r.Replace(v), 64)
		if err != nil {
			log.Errorf("was not able to parse %T to float64!", v)
			return 0
		}

		return float64(result)
	default:
		log.Errorf("conversion to float64 from %T is not supported", v)
		return 0
	}
}

func toSeconds(days float64) float64 {
	return days * 86400
}

func stringToFloat64(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}

	return f
}
