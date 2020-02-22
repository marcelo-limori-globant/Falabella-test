package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Field types. Used for both validation of a type identifier, and field value itself.
var fieldTypeValidators = map[string]string{
	"A": "^[0-9A-Z]+$",
	"N": "^[0-9]+$",
}

func getFieldType(data []byte, cursor int) (string, error) {
	if cursor >= len(data) { // This would have be caught already, but just in case
		return "", errors.New("Index out of range.")
	}

	result := string(data[cursor])
	if _, ok := fieldTypeValidators[result]; !ok {
		return "", errors.New(fmt.Sprintf("Unknown field type %s.", result))
	}
	return result, nil
}

// Get a number on the range 01..99
func getPositiveNumber(data []byte, cursor int) (int, error) {
	if cursor >= len(data)-1 {
		return 0, errors.New("Index out of range.")
	}

	aux := string(data[cursor : cursor+2])
	result, err := strconv.Atoi(aux)
	if result <= 0 || err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid number %s.", aux))
	}
	return result, nil
}

func getFieldNumber(data []byte, cursor int) (string, error) {
	// As of now, a field number is a normal 2 digit number greater than 0.
	// No other validations are required.
	nmbr, err := getPositiveNumber(data, cursor)
	return fmt.Sprintf("%02d", nmbr), err
}

// Validate a field value based on expected type and its regex.
func checkFieldValue(value string, expectedType string) error {
	if result, err := regexp.MatchString(fieldTypeValidators[expectedType], value); err != nil || !result {
		return errors.New(fmt.Sprintf("Invalid field value %s.", value))
	}
	return nil
}

func getFieldValue(data []byte, fieldType string, fieldLength int, cursor int) (string, error) {
	if cursor > len(data)-fieldLength {
		return "", errors.New("Length of field value is out of range.")
	}

	result := string(data[cursor : cursor+fieldLength])

	return result, checkFieldValue(result, fieldType)
}

// Intended for slices. DO NOT pass arrays.
func decodeTLV(data []byte) (map[string]string, error) {
	result := map[string]string{}

	if len(data) == 0 {
		return result, errors.New("Data length is 0.")
	}

	// Ugly, but simplifies flow with the "value, err :=" pattern.
	cursor := 0
	var fieldType string
	var fieldNumber string
	var fieldLength int
	var fieldValue string
	var err error

	for cursor < len(data) {

		if fieldType, err = getFieldType(data, cursor); err != nil {
			return result, err
		}
		cursor += 1

		if fieldNumber, err = getFieldNumber(data, cursor); err != nil {
			return result, err
		}
		cursor += 2

		if fieldLength, err = getPositiveNumber(data, cursor); err != nil {
			return result, err
		}
		cursor += 2

		if fieldValue, err = getFieldValue(data, fieldType, fieldLength, cursor); err != nil {
			return result, err
		}
		cursor += fieldLength

		result[fieldNumber] = fieldValue
	}

	return result, nil
}

func doDecode(data []byte) {
	// Note that this is ok as long as data is a slice and not an array)
	result, err := decodeTLV(data)

	// Note that this will make the test coverage go from 91% to 81%
	fmt.Printf("Data: %s\n\n", string(data))
	fmt.Printf("Error: %v\n\n", err)
	fmt.Printf("Map content:\n\n")
	for key, value := range result {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Printf("\n\n")
}

func main() {
	doDecode([]byte("A0511AB398765UJ1N230200"))
	doDecode([]byte("A0511AB398765UJ1J230200"))
}
