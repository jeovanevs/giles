package archiver

import (
	"strings"
	"time"
)

// remove trailing commas, replace all / with .
func cleantagstring(inp string) string {
	tmp := strings.TrimSuffix(inp, ",")
	tmp = strings.Replace(tmp, "/", ".", -1)
	return tmp
}

// Go doesn't provide min for uint32
func min(a, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Go doesn't provide min for uint32
func min64(a, b uint64) uint64 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Calls function f and then pauses for [pause]. Loops forever.
func periodicCall(pause time.Duration, f func()) {
	for {
		f()
		time.Sleep(pause)
	}
}

// Given a forward-slash delimited path, returns a slice of prefixes, e.g.:
// input: /a/b/c/d
// output: ['/', '/a','/a/b','/a/b/c']
func getPrefixes(s string) []string {
	ret := []string{"/"}
	root := ""
	s = "/" + s
	for _, prefix := range strings.Split(s, "/") {
		if len(prefix) > 0 { //skip empty strings created by Split
			root += "/" + prefix
			ret = append(ret, root)
		}
	}
	if len(ret) > 1 {
		return ret[:len(ret)-1]
	}
	return ret
}

// returns true if the two slices are equal
func isStringSliceEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for idx := range x {
		if x[idx] != y[idx] {
			return false
		}
	}
	return true
}

// unescapes the HTML code for "="
func unescape(s string) string {
	return strings.Replace(s, "%3D", "=", -1)
}

// returns True if the sMAP message contains a []byte reading rather
// than a numerical one
func isObjReading(msg *SmapMessage) bool {
	if msg.Readings == nil || len(msg.Readings) == 0 {
		return false
	} // return early
	_, is_uint := msg.Readings[0][1].(uint64)
	_, is_float := msg.Readings[0][1].(float64)
	if is_uint || is_float {
		return false
	}
	return true
}
