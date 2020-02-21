package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllTypes(t *testing.T) {
	data := []byte("N010212A0202A3")
	result, err := decodeTLV(data)
	assert.Nil(t, err)
	require.NotNil(t, result)
	assertMapValue(t, result, "01", "12")
	assertMapValue(t, result, "02", "A3")
}

func TestZeroLength(t *testing.T) {
	data := []byte("")
	result, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Data length is 0.", err.Error())
	require.NotNil(t, result)
	assert.Empty(t, result)
}

func TestInvalidType(t *testing.T) {
	data := []byte("A010212B010212A010212")
	result, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Unknown field type B.", err.Error())
	require.NotNil(t, result)
	assertMapValue(t, result, "01", "12")
}

func TestTypeMissmatch(t *testing.T) {
	data := []byte("N0102A2B010212A010212")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Invalid field value A2.", err.Error())
}

func TestInvalidFieldNumber(t *testing.T) {
	data := []byte("A3A0212A010212A010212")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Invalid number 3A.", err.Error())
}

func TestInvalidFieldNumberByLength(t *testing.T) {
	data := []byte("A1")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Index out of range.", err.Error())
}

func TestInvalidFieldNumberByRange(t *testing.T) {
	data := []byte("A-1")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Invalid number -1.", err.Error())
}

func TestInvalidLengthByRange(t *testing.T) {
	data := []byte("A01-3123")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Invalid number -3.", err.Error())
}

func TestInvalidLengthOutOfRange(t *testing.T) {
	data := []byte("A0199123")
	_, err := decodeTLV(data)
	require.NotNil(t, err)
	assert.Equal(t, "Length of field value is out of range.", err.Error())
}

func assertMapValue(t *testing.T, theMap map[string]string, key string, value string) {
	assert.Contains(t, theMap, key)
	assert.Equal(t, theMap[key], value)
}
