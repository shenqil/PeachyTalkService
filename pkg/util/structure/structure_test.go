package structure

import "testing"

type Xs struct {
	X int
	Y int
}

type aS struct {
	A   string `json:"a"`
	B   string
	C   string
	SiS string `json:"si_s"`
	Xs
}

type bS struct {
	A   string
	B   string
	D   string
	SiS string `json:"si_s"`
	Xs
}

func TestCopy(t *testing.T) {
	a := &aS{
		A:   "1",
		B:   "2",
		C:   "3",
		SiS: "1223",
		Xs:  Xs{X: 1, Y: 2},
	}

	b := new(bS)
	err := Copy2(a, b)
	if err != nil {
		t.Fail()
		t.Error(err.Error())
	}
	t.Log(b)
}
