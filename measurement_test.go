package main

import "fmt"
import "strconv"
import "testing"

func TestParseMeasurements(t *testing.T) {
	{
		testInput := ""
		_, err := ParseMeasurements(testInput)
		if err == nil {
			t.Errorf("ParseMeasurements(\"%s\") failed, didn't get an error.", testInput)
		} else if _, ok := err.(*strconv.NumError); !ok {
			t.Errorf("ParseMeasurements(\"%s\") failed, expected error: strconv.NumError, got %v", testInput, err)
		}
	}

	{
		testInput := "- 30.23 -"
		testOutput := "30.230 - -"
		m, err := ParseMeasurements(testInput)
		if err != nil {
			t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
		}
		if m.String() != testOutput {
			t.Errorf("ParseMeasurements(\"%s\") failed, expected %v, got %v", testInput, testOutput, m.String())
		}
	}

	{
		testInput := "40 30.23 35.0001"
		testOutput := "30.230 35.000 40.000"
		m, err := ParseMeasurements(testInput)
		if err != nil {
			t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
		}
		if m.String() != testOutput {
			t.Errorf("ParseMeasurements(\"%s\") failed, expected %v, got %v", testInput, testOutput, m.String())
		}
	}
}

func TestLen(t *testing.T) {
	testInput := "50.2 30.1 - 40"
	m, err := ParseMeasurements(testInput)
	if err != nil {
		t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
	}
	if m.Len() != 4 {
		t.Errorf("ParseMeasurements(\"%s\").Len() failed, expected: 4, got: %d", testInput, m.Len())
	}
}

func TestGetSentCount(t *testing.T) {
	testInput := "50.2 30.1 - 40"
	testOutput := 4
	m, err := ParseMeasurements(testInput)
	if err != nil {
		t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
	}
	if m.GetSentCount() != testOutput {
		t.Errorf("ParseMeasurements(\"%s\").GetSentCount() failed, expected: %d, got: %d", testInput, testOutput, m.Len())
	}
}

func TestGetLostCount(t *testing.T) {
	testInput := "50.2 30.1 - 40"
	testOutput := 1
	m, err := ParseMeasurements(testInput)
	if err != nil {
		t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
	}
	if m.GetLostCount() != testOutput {
		t.Errorf("ParseMeasurements(\"%s\").GetLostCount() failed, expected: %d, got: %d", testInput, testOutput, m.Len())
	}
}

func TestGetRTTSum(t *testing.T) {
	testInput := "50.2 30.1 - 40"
	testOutput := "0.120"
	m, err := ParseMeasurements(testInput)
	if err != nil {
		t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
	}
	if fmt.Sprintf("%.3f", m.GetRTTSum()) != testOutput {
		t.Errorf("ParseMeasurements(\"%s\").GetRTTSum() failed, expected: %s, got: %.3f", testInput, testOutput, m.GetRTTSum())
	}
}

func TestGetRTTCount(t *testing.T) {
	testInput := "50.2 30.1 - 40"
	testOutput := 3
	m, err := ParseMeasurements(testInput)
	if err != nil {
		t.Errorf("ParseMeasurements(\"%s\") failed, got error: %v", testInput, err)
	}
	if m.GetRTTCount() != testOutput {
		t.Errorf("ParseMeasurements(\"%s\").GetRTTCount() failed, expected: %d, got: %d", testInput, testOutput, m.GetRTTCount())
	}
}
