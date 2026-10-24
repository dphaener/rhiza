// Package canon reads and writes Genesis's canonical JSON: sorted keys,
// two-space indentation, integers only, minimal string escaping, raw UTF-8
// and exactly one final newline. The same value always has the same bytes,
// so canonical documents can be compared, hashed and stored byte for byte.
//
// The emitter is hand-written rather than built on encoding/json, whose
// encoder always escapes U+2028 and U+2029, HTML-escapes by default and
// silently replaces invalid UTF-8.
package canon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// ErrInvalid is wrapped by every error for input that is not JSON Genesis
// accepts: malformed JSON, duplicate keys, invalid UTF-8, non-integer or
// out-of-range numbers, or nesting deeper than the limit.
var ErrInvalid = errors.New("canon: invalid JSON")

// ErrNotCanonical is wrapped by Check's error when the input is acceptable
// JSON but not in canonical form.
var ErrNotCanonical = errors.New("canon: not canonical")

// maxDepth bounds array and object nesting so hostile input cannot exhaust
// the stack of the recursive parser.
const maxDepth = 1000

// Canonicalize parses exactly one JSON value, optionally surrounded by JSON
// whitespace, and returns its canonical form.
func Canonicalize(data []byte) ([]byte, error) {
	p := parser{data: data}
	p.skipSpace()
	v, err := p.value(0)
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if p.pos != len(p.data) {
		return nil, p.errorf("trailing data")
	}
	var buf bytes.Buffer
	emit(&buf, v, 0)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// Check returns nil only if data is already exactly canonical. Otherwise the
// error wraps ErrInvalid if data does not parse, else ErrNotCanonical.
func Check(data []byte) error {
	out, err := Canonicalize(data)
	if err != nil {
		return err
	}
	if !bytes.Equal(out, data) {
		return ErrNotCanonical
	}
	return nil
}

// Marshal encodes a Go value with encoding/json and returns its canonical
// form. Values that encode to floats or out-of-range integers fail with
// ErrInvalid. Invalid UTF-8 in Go strings is replaced by encoding/json
// before canonicalisation, as it always is.
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return Canonicalize(buf.Bytes())
}

type kind int

const (
	kNull kind = iota
	kTrue
	kFalse
	kInt
	kString
	kArray
	kObject
)

type value struct {
	kind    kind
	num     int64
	str     string
	elems   []value
	members []member
}

type member struct {
	key string
	val value
}

type parser struct {
	data []byte
	pos  int
}

func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("%w: offset %d: %s", ErrInvalid, p.pos, fmt.Sprintf(format, args...))
}

func (p *parser) skipSpace() {
	for p.pos < len(p.data) {
		switch p.data[p.pos] {
		case ' ', '\t', '\n', '\r':
			p.pos++
		default:
			return
		}
	}
}

func (p *parser) value(depth int) (value, error) {
	if p.pos >= len(p.data) {
		return value{}, p.errorf("unexpected end of input")
	}
	switch c := p.data[p.pos]; {
	case c == '{':
		return p.object(depth + 1)
	case c == '[':
		return p.array(depth + 1)
	case c == '"':
		s, err := p.string()
		return value{kind: kString, str: s}, err
	case c == '-' || (c >= '0' && c <= '9'):
		return p.number()
	case c == 't':
		return value{kind: kTrue}, p.literal("true")
	case c == 'f':
		return value{kind: kFalse}, p.literal("false")
	case c == 'n':
		return value{kind: kNull}, p.literal("null")
	default:
		return value{}, p.errorf("unexpected byte %q", c)
	}
}

func (p *parser) literal(word string) error {
	if !bytes.HasPrefix(p.data[p.pos:], []byte(word)) {
		return p.errorf("invalid literal")
	}
	p.pos += len(word)
	return nil
}

func (p *parser) number() (value, error) {
	start := p.pos
	if p.data[p.pos] == '-' {
		p.pos++
	}
	digits := p.pos
	for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
		p.pos++
	}
	switch {
	case p.pos == digits:
		return value{}, p.errorf("number has no digits")
	case p.data[digits] == '0' && p.pos-digits > 1:
		return value{}, p.errorf("number has a leading zero")
	}
	if p.pos < len(p.data) {
		switch p.data[p.pos] {
		case '.', 'e', 'E':
			return value{}, p.errorf("only integers are allowed")
		}
	}
	n, err := strconv.ParseInt(string(p.data[start:p.pos]), 10, 64)
	if err != nil {
		return value{}, p.errorf("integer out of int64 range")
	}
	return value{kind: kInt, num: n}, nil
}

func (p *parser) string() (string, error) {
	p.pos++ // opening quote
	var sb strings.Builder
	for {
		if p.pos >= len(p.data) {
			return "", p.errorf("unterminated string")
		}
		c := p.data[p.pos]
		switch {
		case c == '"':
			p.pos++
			return sb.String(), nil
		case c == '\\':
			if err := p.escape(&sb); err != nil {
				return "", err
			}
		case c < 0x20:
			return "", p.errorf("raw control character in string")
		case c < utf8.RuneSelf:
			sb.WriteByte(c)
			p.pos++
		default:
			r, size := utf8.DecodeRune(p.data[p.pos:])
			if r == utf8.RuneError && size == 1 {
				return "", p.errorf("invalid UTF-8")
			}
			sb.Write(p.data[p.pos : p.pos+size])
			p.pos += size
		}
	}
}

func (p *parser) escape(sb *strings.Builder) error {
	p.pos++ // backslash
	if p.pos >= len(p.data) {
		return p.errorf("unterminated escape")
	}
	c := p.data[p.pos]
	p.pos++
	switch c {
	case '"', '\\', '/':
		sb.WriteByte(c)
	case 'b':
		sb.WriteByte('\b')
	case 'f':
		sb.WriteByte('\f')
	case 'n':
		sb.WriteByte('\n')
	case 'r':
		sb.WriteByte('\r')
	case 't':
		sb.WriteByte('\t')
	case 'u':
		r, err := p.hex4()
		if err != nil {
			return err
		}
		if utf16.IsSurrogate(r) {
			// Only a high surrogate followed by an escaped low surrogate
			// names a code point; anything else is not valid UTF-8.
			if r >= 0xdc00 || !bytes.HasPrefix(p.data[p.pos:], []byte(`\u`)) {
				return p.errorf("lone surrogate")
			}
			p.pos += 2
			lo, err := p.hex4()
			if err != nil {
				return err
			}
			r = utf16.DecodeRune(r, lo)
			if r == utf8.RuneError {
				return p.errorf("lone surrogate")
			}
		}
		sb.WriteRune(r)
	default:
		return p.errorf("invalid escape %q", c)
	}
	return nil
}

func (p *parser) hex4() (rune, error) {
	if len(p.data)-p.pos < 4 {
		return 0, p.errorf("short unicode escape")
	}
	n, err := strconv.ParseUint(string(p.data[p.pos:p.pos+4]), 16, 32)
	if err != nil {
		return 0, p.errorf("invalid unicode escape")
	}
	p.pos += 4
	return rune(n), nil
}

func (p *parser) array(depth int) (value, error) {
	if depth > maxDepth {
		return value{}, p.errorf("nesting deeper than %d", maxDepth)
	}
	p.pos++ // [
	v := value{kind: kArray}
	p.skipSpace()
	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.pos++
		return v, nil
	}
	for {
		p.skipSpace()
		elem, err := p.value(depth)
		if err != nil {
			return value{}, err
		}
		v.elems = append(v.elems, elem)
		p.skipSpace()
		if p.pos >= len(p.data) {
			return value{}, p.errorf("unterminated array")
		}
		switch p.data[p.pos] {
		case ',':
			p.pos++
		case ']':
			p.pos++
			return v, nil
		default:
			return value{}, p.errorf("expected , or ]")
		}
	}
}

func (p *parser) object(depth int) (value, error) {
	if depth > maxDepth {
		return value{}, p.errorf("nesting deeper than %d", maxDepth)
	}
	start := p.pos
	p.pos++ // {
	v := value{kind: kObject}
	p.skipSpace()
	if p.pos < len(p.data) && p.data[p.pos] == '}' {
		p.pos++
		return v, nil
	}
	for {
		p.skipSpace()
		if p.pos >= len(p.data) || p.data[p.pos] != '"' {
			return value{}, p.errorf("expected string key")
		}
		key, err := p.string()
		if err != nil {
			return value{}, err
		}
		p.skipSpace()
		if p.pos >= len(p.data) || p.data[p.pos] != ':' {
			return value{}, p.errorf("expected :")
		}
		p.pos++
		p.skipSpace()
		val, err := p.value(depth)
		if err != nil {
			return value{}, err
		}
		v.members = append(v.members, member{key: key, val: val})
		p.skipSpace()
		if p.pos >= len(p.data) {
			return value{}, p.errorf("unterminated object")
		}
		switch p.data[p.pos] {
		case ',':
			p.pos++
		case '}':
			p.pos++
			// Go string comparison is byte-wise, which is the canonical order.
			slices.SortFunc(v.members, func(a, b member) int { return strings.Compare(a.key, b.key) })
			for i := 1; i < len(v.members); i++ {
				if v.members[i].key == v.members[i-1].key {
					p.pos = start
					return value{}, p.errorf("duplicate key %q", v.members[i].key)
				}
			}
			return v, nil
		default:
			return value{}, p.errorf("expected , or }")
		}
	}
}

func emit(buf *bytes.Buffer, v value, indent int) {
	switch v.kind {
	case kNull:
		buf.WriteString("null")
	case kTrue:
		buf.WriteString("true")
	case kFalse:
		buf.WriteString("false")
	case kInt:
		buf.WriteString(strconv.FormatInt(v.num, 10))
	case kString:
		emitString(buf, v.str)
	case kArray:
		if len(v.elems) == 0 {
			buf.WriteString("[]")
			return
		}
		buf.WriteByte('[')
		for i, e := range v.elems {
			if i > 0 {
				buf.WriteByte(',')
			}
			newline(buf, indent+1)
			emit(buf, e, indent+1)
		}
		newline(buf, indent)
		buf.WriteByte(']')
	case kObject:
		if len(v.members) == 0 {
			buf.WriteString("{}")
			return
		}
		buf.WriteByte('{')
		for i, m := range v.members {
			if i > 0 {
				buf.WriteByte(',')
			}
			newline(buf, indent+1)
			emitString(buf, m.key)
			buf.WriteString(": ")
			emit(buf, m.val, indent+1)
		}
		newline(buf, indent)
		buf.WriteByte('}')
	}
}

func newline(buf *bytes.Buffer, indent int) {
	buf.WriteByte('\n')
	for range indent {
		buf.WriteString("  ")
	}
}

const hexDigits = "0123456789abcdef"

// emitString escapes only what JSON requires; the parser has already
// guaranteed s is valid UTF-8.
func emitString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if c < 0x20 {
				buf.WriteString(`\u00`)
				buf.WriteByte(hexDigits[c>>4])
				buf.WriteByte(hexDigits[c&0xf])
			} else {
				buf.WriteByte(c)
			}
		}
	}
	buf.WriteByte('"')
}
