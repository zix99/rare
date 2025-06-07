package dissect

import (
	"testing"

	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestDissectBasic(t *testing.T) {
	d := MustCompile("%{val};%{};%{?skip} - %{val2}")

	assert.Equal(t, []int{0, 17, 0, 5, 12, 17}, d.FindSubmatchIndex([]byte("Hello;a;b - there")))

	assert.Equal(t, map[string]int{
		"val":  1,
		"val2": 2,
	}, d.SubexpNameTable())
}

func TestUtf8(t *testing.T) {
	d := MustCompile("ûɾ %{key} ḝłįʈ")

	s := []byte("Ḽơᶉëᶆ ȋṕšᶙṁ ḍỡḽǭᵳ ʂǐť ӓṁệẗ, ĉṓɲṩḙċťᶒțûɾ ấɖḯƥĭṩčįɳġ ḝłįʈ, șếᶑ ᶁⱺ ẽḭŭŝḿꝋď ṫĕᶆᶈṓɍ ỉñḉīḑȋᵭṵńť ṷŧ ḹẩḇőꝛế éȶ đꝍꞎôꝛȇ ᵯáꞡᶇā ąⱡîɋṹẵ.")
	m := d.FindSubmatchIndex(s)
	assert.Equal(t, []int{85, 123, 90, 113}, m)
	assert.Equal(t, []byte("ûɾ ấɖḯƥĭṩčįɳġ ḝłįʈ"), s[m[0]:m[1]])
	assert.Equal(t, []byte("ấɖḯƥĭṩčįɳġ"), s[m[2]:m[3]])
}

func TestPrefixOnSkipKey(t *testing.T) {
	d := MustCompile("prefix %{}: %{val}")

	assert.Nil(t, d.FindSubmatchIndex([]byte("a: b")))
	assert.Equal(t, []int{0, 11, 10, 11}, d.FindSubmatchIndex([]byte("prefix a: b")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("Prefix a: b")))
}

func TestEmpty(t *testing.T) {
	d := MustCompile("")
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndex([]byte("hello")))
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndex([]byte("")))
}

func TestNoTokens(t *testing.T) {
	d := MustCompile("test")

	assert.Nil(t, d.FindSubmatchIndex([]byte("hello there")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("test")))
	assert.Equal(t, []int{1, 5}, d.FindSubmatchIndex([]byte("atest")))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndex([]byte("abctest")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("testa")))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndex([]byte("testabc")))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndex([]byte("abctestabc")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("tEst")))
}

func TestPrefix(t *testing.T) {
	d := MustCompile("mid %{val};%{val2} after")

	assert.Equal(t, []int{12, 29, 16, 19, 20, 23}, d.FindSubmatchIndex([]byte("string with mid 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("string with mi 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("string with mid 123;456 boom k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("")))
}

func TestSuffix(t *testing.T) {
	d := MustCompile("%{val};%{val2} after")

	assert.Equal(t, []int{0, 13, 0, 3, 4, 7}, d.FindSubmatchIndex([]byte("123;456 after k")))
	assert.Equal(t, []int{0, 17, 0, 7, 8, 11}, d.FindSubmatchIndex([]byte("hah 123;456 after k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("123;456 boom k")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("")))

	assert.Equal(t, []int{2, 13, 6, 13}, MustCompile("end %{nada}").FindSubmatchIndex([]byte("a end nothing")))
}

func TestNoPrefixSuffix(t *testing.T) {
	d := MustCompile("%{onlymatch}")
	assert.Equal(t, []int{0, 5, 0, 5}, d.FindSubmatchIndex([]byte("a b c")))
}

func TestErrorNew(t *testing.T) {
	// Unclosed
	_, err := Compile("unclosed %{")
	assert.ErrorIs(t, err, ErrorUnclosedToken)

	// Dupe key
	_, err = Compile("a %{a} %{a}")
	assert.ErrorIs(t, err, ErrorKeyConflict)

	// Sequential tokens
	_, err = Compile("a %{a}%{b}")
	assert.ErrorIs(t, err, ErrorSequentialToken)
}

func TestMustPanics(t *testing.T) {
	assert.Panics(t, func() {
		MustCompile("%{bad expr")
	})
}

func TestIgnoreCase(t *testing.T) {
	d, err := CompileEx("TeSt1", true)

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 5}, d.FindSubmatchIndex([]byte("test1")))
	assert.Equal(t, []int{0, 5}, d.FindSubmatchIndex([]byte("tEst1")))
	assert.Equal(t, []int{0, 5}, d.FindSubmatchIndex([]byte("TEST1")))
	assert.Equal(t, []int{1, 6}, d.FindSubmatchIndex([]byte("ATest123")))
	assert.Nil(t, d.FindSubmatchIndex([]byte("asdf")))

	d, err = CompileEx("pref %{val} post", true)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 13, 7, 8}, d.FindSubmatchIndex([]byte("a Pref 5 pOst")))
}

func TestMemoryTarget(t *testing.T) {
	d := MustCompile("prefix %{}: %{val}")
	str := []byte("prefix a: 123")

	buf := make([]int, 0, d.MatchBufSize())
	ret := d.FindSubmatchIndexDst([]byte(str), buf)
	assert.Equal(t, []int{0, 13, 10, 13}, ret)
	testutil.AssertSameMemory(t, buf, ret)

	// undersized
	buf = make([]int, 0, 2)
	ret = d.FindSubmatchIndexDst([]byte(str), buf)
	assert.Equal(t, []int{0, 13, 10, 13}, ret)
	testutil.AssertNotSameMemory(t, buf, ret)
}

func TestDstZeroAlloc(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkDissectDst)
}

// BenchmarkDissect-4   	13347456	        86.07 ns/op	      32 B/op	       0 allocs/op
func BenchmarkDissect(b *testing.B) {
	d, _ := CompileEx("t%{val} ", false)
	val := []byte("this is a test ")

	for i := 0; i < b.N; i++ {
		d.FindSubmatchIndex(val)
	}
}

// BenchmarkDissectDst-4   	23545326	        45.58 ns/op	       0 B/op	       0 allocs/op
func BenchmarkDissectDst(b *testing.B) {
	d, _ := CompileEx("t%{val} ", false)
	buf := make([]int, 0, d.MatchBufSize())
	val := []byte("this is a test ")

	for i := 0; i < b.N; i++ {
		d.FindSubmatchIndexDst(val, buf)
	}
}
