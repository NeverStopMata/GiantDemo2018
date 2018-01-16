package util

import (
	"io/ioutil"
	"strings"
)

type WordFilter struct {
	words     map[string]struct{}
	firstChar map[rune]struct{}
	allChar   map[rune]struct{}
	maxlen    int
}

var dirtyfilter *WordFilter

func LoadDirtyWords(file string) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	words := strings.Split(string(d), "|")

	dirtyfilter = NewWordFilter()
	for _, val := range words {
		dirtyfilter.Insert(val)
	}
	return nil
}

func IsDirty(str string) bool {
	if dirtyfilter != nil {
		return dirtyfilter.IsDirty(str)
	}
	return false
}

func NewWordFilter() (filter *WordFilter) {
	filter = &WordFilter{}
	filter.words = make(map[string]struct{})
	filter.firstChar = make(map[rune]struct{})
	filter.allChar = make(map[rune]struct{})
	return
}

func (this *WordFilter) Insert(str string) {
	if len(str) == 0 {
		return
	}
	rstr := []rune(str)
	this.words[str] = struct{}{}
	this.firstChar[rstr[0]] = struct{}{}
	if len(rstr) > this.maxlen {
		this.maxlen = len(rstr)
	}
	for _, w := range rstr {
		this.allChar[w] = struct{}{}
	}
}

func (this *WordFilter) IsDirty(str string) bool {
	wstr := []rune(str)
	if len(wstr) == 0 {
		return false
	}
	index, offset, mlen := 0, 0, 0
	for index < len(wstr) {
		if _, ok := this.firstChar[wstr[index]]; !ok {
			index++
			continue
		}
		mlen = len(wstr) - index
		if mlen > this.maxlen {
			mlen = this.maxlen
		}
		for j := 1; j <= mlen; j++ {
			offset = index + j
			if _, ok := this.allChar[wstr[offset-1]]; !ok {
				break
			}
			if _, ok := this.words[string(wstr[index:offset])]; ok {
				return true
			}
		}
		index++
	}
	return false
}
