package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dustin/go-jsonpointer"
)

type Condition struct {
	Field string
	Match string
	re    *regexp.Regexp `yaml:"-"`
}

func (c *Condition) Check(info *Nodeinfo) bool {
	val := jsonpointer.Reflect(info, c.Field)
	if val == nil {
		return false
	}

	strval := strings.ToLower(fmt.Sprint(val))
	return c.re.MatchString(strval)
}
