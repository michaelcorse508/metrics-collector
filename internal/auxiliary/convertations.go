package auxiliary

import "strconv"

func ConvertStringToFloat64(stringValue string) (*float64, error) {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return nil, err
	}

	outValue := new(float64)
	*outValue = value

	return outValue, nil
}

func ConvertStringToInt64(stringValue string) (*int64, error) {
	value, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return nil, err
	}

	outValue := new(int64)
	*outValue = value

	return outValue, nil
}
