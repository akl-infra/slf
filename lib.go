package slf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Finger uint8
type KeymeowFinger Finger

const (
	LP Finger = iota
	LR
	LM
	LI
	LT
	RT
	RI
	RM
	RR
	RP
)

var (
	FingerName = []string{
		"LP",
		"LR",
		"LM",
		"LI",
		"LT",
		"RT",
		"RI",
		"RM",
		"RR",
		"RP",
	}
	FingerValue = map[string]Finger{
		"LP": LP,
		"LR": LR,
		"LM": LM,
		"LI": LI,
		"LT": LT,
		"RT": RT,
		"RI": RI,
		"RM": RM,
		"RR": RR,
		"RP": RP,
	}
)

func (f Finger) String() string {
	return FingerName[f]
}

func ParseFinger(s string) (Finger, error) {
	val, ok := FingerValue[s]
	if !ok {
		return Finger(0), fmt.Errorf("%s is not a valid finger", s)
	}
	return val, nil
}

func (f KeymeowFinger) MarshalJSON() ([]byte, error) {
	return json.Marshal(Finger(f).String())
}

func (f *Finger) UnmarshalJSON(data []byte) (err error) {
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*f = Finger(num)
		return nil
	}
	var finger string
	if err := json.Unmarshal(data, &finger); err != nil {
		return err
	}
	if *f, err = ParseFinger(finger); err != nil {
		return err
	}
	return nil
}

func (f *KeymeowFinger) UnmarshalJSON(data []byte) (err error) {
	var finger string
	if err := json.Unmarshal(data, &finger); err != nil {
		return err
	}
	var of Finger
	if of, err = ParseFinger(finger); err != nil {
		return err
	}
	*f = KeymeowFinger(of)
	return nil
}

type Key struct {
	Char   string `json:"char"`
	Row    uint8  `json:"row"`
	Col    uint8  `json:"col"`
	Finger Finger `json:"finger"`
}

type Layout struct {
	Name     string    `json:"name"`
	Owner    int64     `json:"owner"`
	Author   string    `json:"author"`
	Link     string    `json:"link"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Boards   []string  `json:"boards"`
	Keys     []Key     `json:"keys"`
}

// Used for formats that store layouts as a 2d matrix
type MatrixKey struct {
	Char   rune
	Finger Finger
}

type KeymeowComponent struct {
	Finger []KeymeowFinger
	Keys   []string
}

type KeymeowLayout struct {
	Name       string
	Authors    []string
	Components []KeymeowComponent
}

type ConversionError struct {
	Err error
}

func (e *ConversionError) Error() string {
	return e.Err.Error()
}

func (l *Layout) Matrix() [][]MatrixKey {
	rows := make([][]MatrixKey, 0)
	for _, key := range l.Keys {
		for len(rows) <= int(key.Row) {
			// add necessary rows
			rows = append(rows, make([]MatrixKey, 0, 20))
		}
		if len(rows[key.Row]) <= int(key.Col) {
			// expand row to fit key
			rows[key.Row] = rows[key.Row][:key.Col+1]
		}
		rows[key.Row][key.Col] = MatrixKey{[]rune(key.Char)[0], key.Finger}
	}
	return rows
}

func (l *Layout) ToGenkey() (string, error) {
	matrix := l.Matrix()
	var b strings.Builder
	if len(matrix) != 3 {
		return "", &ConversionError{errors.New("genkey only supports layouts with 3 rows")}
	}
	b.WriteString(l.Name)
	b.WriteRune('\n')
	for _, fingermap := range []bool{false, true} {
		for _, row := range matrix {
			for _, key := range row {
				if fingermap {
					f := key.Finger
					if f == LT || f == RT {
						return "", &ConversionError{errors.New("genkey does not support thumbkeys")}
					}
					if f >= RI {
						f -= 2
					}
					b.WriteString(strconv.Itoa(int(f)))
				} else {
					b.WriteRune(key.Char)
				}
				b.WriteRune(' ')
			}
			b.WriteRune('\n')
		}
	}
	return b.String(), nil
}

func (l *Layout) ToOxeylyzer() (string, error) {
	matrix := l.Matrix()
	var b strings.Builder
	if len(matrix) != 3 {
		return "", &ConversionError{errors.New("oxeylyzer only supports 3x10 layouts")}
	}
	for _, row := range matrix {
		if len(row) != 10 {
			return "", &ConversionError{errors.New("oxeylyzer only supports 3x10 layouts")}
		}
		for _, key := range row {
			b.WriteRune(key.Char)
			b.WriteRune(' ')
		}
		b.WriteRune('\n')
	}
	return b.String(), nil
}

func (l Layout) ToKeymeow() KeymeowLayout {
	var keymeow KeymeowLayout
	keymeow.Name = l.Name
	keymeow.Authors = []string{l.Author}
	keymeow.Components = make([]KeymeowComponent, 10)
	for i := range keymeow.Components {
		finger := []KeymeowFinger{KeymeowFinger(i)}
		keymeow.Components[i] = KeymeowComponent{finger, make([]string, 0, 12)}
	}
	slices.SortFunc(l.Keys, func(a, b Key) int {
		if a.Col < b.Col {
			return -1
		} else if b.Col < a.Col {
			return 1
		} else {
			if a.Row < b.Row {
				return -1
			}
			return 1
		}
	})
	for _, key := range l.Keys {
		keys := &keymeow.Components[key.Finger].Keys
		*keys = append(*keys, key.Char)
	}
	return keymeow
}

func ReadLayoutFile(filename string) (Layout, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return Layout{}, err
	}
	l, err := ParseLayout(b)
	return l, err
}

func ParseLayout(b []byte) (Layout, error) {
	var l Layout
	err := json.Unmarshal(b, &l)
	return l, err
}
