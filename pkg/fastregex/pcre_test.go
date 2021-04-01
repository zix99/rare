package fastregex

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileSuccess(t *testing.T) {
	_, err := Compile("test.*")
	assert.NoError(t, err)
}

func TestCompileError(t *testing.T) {
	_, err := Compile("test(.*")
	assert.Error(t, err)
}

func TestCaptureGroupCount(t *testing.T) {
	re, err := Compile("this is (.+) test (.+)")
	assert.NoError(t, err)
	assert.Equal(t, 2, re.GroupCount())
}

func TestMatch(t *testing.T) {
	assert.True(t, MustCompile("test").Match([]byte("this is a test")))
	assert.False(t, MustCompile("test").Match([]byte("this is a tes")))
}

func TestSubMatch(t *testing.T) {
	re := MustCompile("a (\\w+)")
	ret := re.FindSubmatchIndex([]byte("this is a test"))
	assert.Len(t, ret, 4)
	assert.Equal(t, []int{8, 14, 10, 14}, ret)
}

func BenchmarkPCREMatch(b *testing.B) {
	re := MustCompile("t(\\w+)")
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

func BenchmarkRE2Match(b *testing.B) {
	re := regexp.MustCompile("t(\\w+)")
	for i := 0; i < b.N; i++ {
		re.MatchString("this is a test")
	}
}

func BenchmarkPCRESubMatch(b *testing.B) {
	re := MustCompile("t(\\w+)")
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}

func BenchmarkRE2SubMatch(b *testing.B) {
	re := regexp.MustCompile("t(\\w+)")
	for i := 0; i < b.N; i++ {
		re.FindSubmatchIndex([]byte("this is a test"))
	}
}
