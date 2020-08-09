package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Measurements struct {
	rtt  []float64
	lost []bool
}

func ParseMeasurements(text string) (Measurements, error) {
	valuesText := strings.Split(text, " ")
	m := Measurements{
		rtt:  make([]float64, len(valuesText)),
		lost: make([]bool, len(valuesText)),
	}
	for i, valText := range valuesText {
		if valText == "-" {
			m.lost[i] = true
		} else {
			m.lost[i] = false
			rtt, err := strconv.ParseFloat(valText, 32)
			if err != nil {
				return m, err
			}
			m.rtt[i] = float64(rtt) / 1000.0
		}
	}
	sort.Sort(m)
	return m, nil
}

func (m Measurements) Len() int {
	return len(m.rtt)
}

func (m Measurements) Less(i, j int) bool {
	if m.lost[i] {
		return false
	}
	if m.lost[j] {
		return true
	}
	return m.rtt[i] < m.rtt[j]
}

func (m Measurements) Swap(i, j int) {
	m.rtt[i], m.rtt[j] = m.rtt[j], m.rtt[i]
	m.lost[i], m.lost[j] = m.lost[j], m.lost[i]
}

func (m Measurements) String() string {
	var str strings.Builder
	for i := range m.rtt {
		if i != 0 {
			str.WriteString(" ")
		}
		if m.lost[i] {
			str.WriteString("-")
		} else {
			str.WriteString(fmt.Sprintf("%.3f", m.rtt[i] * 1000.0))
		}
	}
	return str.String()
}

func (m Measurements) GetSentCount() int {
	return len(m.lost)
}

func (m Measurements) GetLostCount() int {
	var count int
	for _, val := range m.lost {
		if val {
			count++
		}
	}
	return count
}

func (m Measurements) GetRTTSum() float64 {
	var sum float64
	for i, val := range m.rtt {
		if !m.lost[i] {
			sum += val
		}
	}
	return sum
}

func (m Measurements) GetRTTCount() int {
	var count int
	for _, val := range m.lost {
		if !val {
			count++
		}
	}
	return count
}
