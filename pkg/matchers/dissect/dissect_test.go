package dissect

import (
	"rare/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDissectBasic(t *testing.T) {
	d := MustCompile("%{val};%{};%{?skip} - %{val2}").CreateInstance()

	assert.Equal(t, []int{0, 17, 0, 5, 12, 17}, d.FindSubmatchIndexDst([]byte("Hello;a;b - there"), nil))

	assert.Equal(t, 6, d.MatchBufSize())
	assert.Equal(t, map[string]int{
		"val":  1,
		"val2": 2,
	}, d.SubexpNameTable())
}
func TestUtf8(t *testing.T) {
	d := MustCompile("ûɾ %{key} ḝłįʈ").CreateInstance()
	assert.Equal(t, 4, d.MatchBufSize())

	s := []byte("Ḽơᶉëᶆ ȋṕšᶙṁ ḍỡḽǭᵳ ʂǐť ӓṁệẗ, ĉṓɲṩḙċťᶒțûɾ ấɖḯƥĭṩčįɳġ ḝłįʈ, șếᶑ ᶁⱺ ẽḭŭŝḿꝋď ṫĕᶆᶈṓɍ ỉñḉīḑȋᵭṵńť ṷŧ ḹẩḇőꝛế éȶ đꝍꞎôꝛȇ ᵯáꞡᶇā ąⱡîɋṹẵ.")
	m := d.FindSubmatchIndexDst(s, nil)
	assert.Equal(t, []int{85, 123, 90, 113}, m)
	assert.Equal(t, []byte("ûɾ ấɖḯƥĭṩčįɳġ ḝłįʈ"), s[m[0]:m[1]])
	assert.Equal(t, []byte("ấɖḯƥĭṩčįɳġ"), s[m[2]:m[3]])
}

func TestPrefixOnSkipKey(t *testing.T) {
	d := MustCompile("prefix %{}: %{val}").CreateInstance()

	assert.Nil(t, d.FindSubmatchIndexDst([]byte("a: b"), nil))
	assert.Equal(t, []int{0, 11, 10, 11}, d.FindSubmatchIndexDst([]byte("prefix a: b"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte("Prefix a: b"), nil))
}

func TestEmpty(t *testing.T) {
	d := MustCompile("").CreateInstance()
	assert.Equal(t, 2, d.MatchBufSize())
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndexDst([]byte("hello"), nil))
	assert.Equal(t, []int{0, 0}, d.FindSubmatchIndexDst([]byte(""), nil))
}

func TestNoTokens(t *testing.T) {
	d := MustCompile("test").CreateInstance()
	assert.Equal(t, 2, d.MatchBufSize())

	assert.Nil(t, d.FindSubmatchIndexDst([]byte("hello there"), nil))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndexDst([]byte("test"), nil))
	assert.Equal(t, []int{1, 5}, d.FindSubmatchIndexDst([]byte("atest"), nil))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndexDst([]byte("abctest"), nil))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndexDst([]byte("testa"), nil))
	assert.Equal(t, []int{0, 4}, d.FindSubmatchIndexDst([]byte("testabc"), nil))
	assert.Equal(t, []int{3, 7}, d.FindSubmatchIndexDst([]byte("abctestabc"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte("tEst"), nil))
}

func TestPrefix(t *testing.T) {
	d := MustCompile("mid %{val};%{val2} after").CreateInstance()

	assert.Equal(t, []int{12, 29, 16, 19, 20, 23}, d.FindSubmatchIndexDst([]byte("string with mid 123;456 after k"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte("string with mi 123;456 after k"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte("string with mid 123;456 boom k"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte(""), nil))
}

func TestSuffix(t *testing.T) {
	d := MustCompile("%{val};%{val2} after").CreateInstance()

	assert.Equal(t, []int{0, 13, 0, 3, 4, 7}, d.FindSubmatchIndexDst([]byte("123;456 after k"), nil))
	assert.Equal(t, []int{0, 17, 0, 7, 8, 11}, d.FindSubmatchIndexDst([]byte("hah 123;456 after k"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte("123;456 boom k"), nil))
	assert.Nil(t, d.FindSubmatchIndexDst([]byte(""), nil))

	assert.Equal(t, []int{2, 13, 6, 13}, MustCompile("end %{nada}").CreateInstance().FindSubmatchIndexDst([]byte("a end nothing"), nil))
}

func TestNoPrefixSuffix(t *testing.T) {
	d := MustCompile("%{onlymatch}").CreateInstance()
	assert.Equal(t, 4, d.MatchBufSize())
	assert.Equal(t, []int{0, 5, 0, 5}, d.FindSubmatchIndexDst([]byte("a b c"), nil))
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
	assert.Equal(t, []int{0, 5}, d.CreateInstance().FindSubmatchIndexDst([]byte("test1"), nil))
	assert.Equal(t, []int{0, 5}, d.CreateInstance().FindSubmatchIndexDst([]byte("tEst1"), nil))
	assert.Equal(t, []int{0, 5}, d.CreateInstance().FindSubmatchIndexDst([]byte("TEST1"), nil))
	assert.Equal(t, []int{1, 6}, d.CreateInstance().FindSubmatchIndexDst([]byte("ATest123"), nil))
	assert.Nil(t, d.CreateInstance().FindSubmatchIndexDst([]byte("asdf"), nil))

	d, err = CompileEx("pref %{val} post", true)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 13, 7, 8}, d.CreateInstance().FindSubmatchIndexDst([]byte("a Pref 5 pOst"), nil))
}

func TestSameMemory(t *testing.T) {
	d := MustCompile("test").CreateInstance()
	assert.Equal(t, 2, d.MatchBufSize())
	buf := make([]int, 0, 2)

	ret := d.FindSubmatchIndexDst([]byte("test"), buf)
	assert.Equal(t, []int{0, 4}, ret)
	testutil.AssertSameMemory(t, buf, ret)
}

func TestZeroAlloc(t *testing.T) {
	testutil.AssertZeroAlloc(t, BenchmarkDissect)
}

// BenchmarkDissect-4   	24508795	        44.94 ns/op	       0 B/op	       0 allocs/op
func BenchmarkDissect(b *testing.B) {
	d, _ := CompileEx("t%{val} ", false)
	di := d.CreateInstance()
	val := []byte("this is a test ")
	buf := make([]int, 0, d.MatchBufSize())

	for i := 0; i < b.N; i++ {
		di.FindSubmatchIndexDst(val, buf)
	}
}
