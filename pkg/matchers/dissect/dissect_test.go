package dissect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDissectBasic(t *testing.T) {
	d := MustNew("%{val};%{};%{?skip} - %{val2}").CreateInstance()

	assert.Equal(t, []int{0, 17, 0, 5, 12, 17}, d.FindSubmatchIndex([]byte("Hello;a;b - there")))

	assert.Equal(t, map[string]int{
		"val":  1,
		"val2": 2,
	}, d.SubexpNameTable())
}

func TestEmpty(t *testing.T) {
	d := MustNew("").CreateInstance()
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndex([]byte("hello")))
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndex([]byte("")))
}

func TestNoTokens(t *testing.T) {
	d := MustNew("test").CreateInstance()

	assert.Nil(t, d.FindSubmatchIndex([]byte("hello there")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("test")))
	assert.Equal(t, []int{1, 5}, d.FindSubmatchIndex([]byte("atest")))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndex([]byte("abctest")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("testa")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("testabc")))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndex([]byte("abctestabc")))
}

func TestPrefix(t *testing.T) {
	d := MustNew("mid %{val};%{val2} after").CreateInstance()

	assert.Equal(t, []int{12, 29, 16, 19, 20, 23}, d.FindSubmatchIndex([]byte("string with mid 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("string with mi 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("string with mid 123;456 boom k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("")))
}

func TestSuffix(t *testing.T) {
	d := MustNew("%{val};%{val2} after").CreateInstance()

	assert.Equal(t, []int{0, 13, 0, 3, 4, 7}, d.FindSubmatchIndex([]byte("123;456 after k")))
	assert.Equal(t, []int{0, 17, 0, 7, 8, 11}, d.FindSubmatchIndex([]byte("hah 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("123;456 boom k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("")))

	assert.Equal(t, []int{2, 13, 6, 13}, MustNew("end %{nada}").CreateInstance().FindSubmatchIndex([]byte("a end nothing")))
}

func TestNoPrefixSuffix(t *testing.T) {
	d := MustNew("%{onlymatch}").CreateInstance()
	assert.Equal(t, []int{0, 5, 0, 5}, d.FindSubmatchIndex([]byte("a b c")))
}

func TestErrorNew(t *testing.T) {
	// Unclosed
	_, err := New("unclosed %{")
	assert.ErrorIs(t, err, ErrorUnclosedToken)

	// Dupe key
	_, err = New("a %{a} %{a}")
	assert.ErrorIs(t, err, ErrorKeyConflict)

	// Sequential tokens
	_, err = New("a %{a}%{b}")
	assert.ErrorIs(t, err, ErrorSequentialToken)
}

func TestMustPanics(t *testing.T) {
	assert.Panics(t, func() {
		MustNew("%{bad expr")
	})
}

// 88 ns
func BenchmarkDissect(b *testing.B) {
	d, _ := New("t%{val} ")
	di := d.CreateInstance()
	val := []byte("this is a test ")

	for i := 0; i < b.N; i++ {
		di.FindSubmatchIndex(val)
	}
}
