package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	// VALUE_TYPE_STRING 字符串类型配置
	VALUE_TYPE_STRING = iota
	// VALUE_TYPE_SLICE 数组类型配置
	VALUE_TYPE_SLICE
	// VALUE_TYPE_MAP MAP类型配置
	VALUE_TYPE_MAP
)

var (
	// ErrorSectionSyntax _
	ErrorSectionSyntax = fmt.Errorf("section syntax error")
	// ErrorKeySyntax _
	ErrorKeySyntax = fmt.Errorf("key syntax error")
	// ErrorValueSyntax _
	ErrorValueSyntax = fmt.Errorf("value syntax error")
)

var (
	// CommentMark 注释标记
	CommentMark = map[rune]bool{
		'#': true,
		';': true,
	}

	// QuoteMark 引号标记
	QuoteMark = map[rune]bool{
		'"':  true,
		'`':  true,
		'\'': true,
	}
)

type iniParser struct {
	reader *bufio.Reader
	eof    bool
	line   int

	currSection string
	Sections    map[string]map[string]interface{}
}

func parseIni(reader io.Reader) (p *iniParser, err error) {
	p = &iniParser{
		reader: bufio.NewReader(reader),
	}
	err = p.parse()
	return
}

func (p *iniParser) parse() (err error) {

	p.Sections = make(map[string]map[string]interface{})

	for {

		// 读取完毕
		if p.eof {
			break
		}

		var line string
		line, err = p.nextLine()
		if err != nil {
			return
		}

		// 过滤掉边界的spcace
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line[0] {

		// 注释处理
		case ';', '#':
			// skip

		// Section处理
		case '[':
			err = p.parseSection(line)

		// key-value处理
		default:
			err = p.parseKeyValue(line)
		}

		// 错误处理
		if err != nil {
			err = fmt.Errorf(err.Error()+":%s (%d)", line, p.line)
			return
		}

	}

	return
}

// nextLine _
func (p *iniParser) nextLine() (line string, err error) {

	line, err = p.reader.ReadString('\n')

	if line != "" {
		p.line++
	}

	if err == io.EOF {
		err = nil
		p.eof = true
	}

	return
}

func (p *iniParser) parseSection(line string) (err error) {

	rLine := []rune(line)
	sectionDone := false
	section := ""

	for i := 1; i < len(rLine); i++ {
		r := rLine[i]

		// section 结束
		if r == ']' {
			section = string(rLine[1:i])
			sectionDone = true
		}

		// 如果section结束了,后面只允许注释
		if sectionDone {
			if !p.eatComment(i+1, rLine) {
				return ErrorSectionSyntax
			}
			break
		}
	}

	// 过滤边界
	section = strings.TrimSpace(section)

	if !sectionDone || section == "" {
		return ErrorSectionSyntax
	}
	p.currSection = section

	return
}

func (p *iniParser) parseKeyValue(line string) (err error) {

	sp := strings.SplitN(line, "=", 2)
	if len(sp) != 2 {
		err = ErrorKeySyntax
		return
	}

	// Key处理
	keyStr, subKey, valueType, err := p.parseKey(sp[0])
	if err != nil {
		return
	}

	section := p.Sections[p.currSection]
	if section == nil {
		section = make(map[string]interface{})
		p.Sections[p.currSection] = section
	}

	valueStr, err := p.parseValue(sp[1])
	if err != nil {
		return
	}

	// 赋值操作
	switch valueType {

	// String 类型
	case VALUE_TYPE_STRING:
		switch section[keyStr].(type) {
		case string, nil:
			section[keyStr] = valueStr
		default:
			err = ErrorValueSyntax
			return
		}

	// Slice 类型
	case VALUE_TYPE_SLICE:
		switch t := section[keyStr].(type) {
		case []string:
			section[keyStr] = append(t, valueStr)
		case nil:
			section[keyStr] = []string{valueStr}
		default:
			err = ErrorValueSyntax
			return
		}

	// Map 类型
	case VALUE_TYPE_MAP:
		switch t := section[keyStr].(type) {
		case map[string]string:
			t[subKey] = valueStr
		case nil:
			section[keyStr] = map[string]string{subKey: valueStr}
		default:
			err = ErrorValueSyntax
			return
		}
	}

	return
}

func (p *iniParser) parseKey(k string) (key, subKey string, valueType int, err error) {

	valueType = VALUE_TYPE_STRING
	k = strings.TrimSpace(k)

	var (
		isSubKey = false
		kBuffer  = bytes.Buffer{}
		skBuffer = bytes.Buffer{}
	)

	for i, r := range []rune(k) {

		if !isSubKey { // 处理key
			switch r {
			case '[':
				isSubKey = true
				valueType = VALUE_TYPE_SLICE
			default:
				kBuffer.WriteRune(r)
			}

		} else { // 处理subkey
			switch r {
			case ']':

				// 如果存在subkey,] 应该是key终结符
				if i != len(k)-1 {
					err = ErrorKeySyntax
					return
				}

			default:

				// 没有遇到终结符号就结束了应该直接报错
				if i == len(k)-1 {
					err = ErrorKeySyntax
					return
				}

				skBuffer.WriteRune(r)
			}
		}

	}

	key = strings.TrimSpace(kBuffer.String())
	subKey = strings.TrimSpace(skBuffer.String())

	if subKey != "" {
		valueType = VALUE_TYPE_MAP
	}

	return
}

func (p *iniParser) parseValue(v string) (value string, err error) {

	v = strings.TrimSpace(v)

	var (
		hasQuote    = false
		quote       = ' '
		valueBuffer = bytes.Buffer{}
		rv          = []rune(v)
	)

	for i, r := range rv {

		// 检查是否启动引号模式,还没有支持引号里面转意
		if i == 0 && QuoteMark[r] {
			quote = r
			hasQuote = true
			continue
		}

		// 引号模式
		if quote != ' ' {

			// 引号模式括号闭合
			if r == quote {
				if !p.eatComment(i+1, []rune(rv)) {
					err = ErrorValueSyntax
					return
				} else {
					break
				}
			}

			// 引号模式到结束括号没有闭合
			if i == len(rv)-1 {
				err = ErrorValueSyntax
				return
			}

			valueBuffer.WriteRune(r)

			// 无引号模式
		} else {
			if CommentMark[r] {
				break
			}
			valueBuffer.WriteRune(r)
		}
	}

	value = valueBuffer.String()

	// 如果带上了括号,那么不处理两边的space,如果没有带括号,需要处理掉
	if !hasQuote {
		value = strings.TrimSpace(value)
	}

	return
}

// 一般用来检查行终结过后,到注释是否正确
func (p *iniParser) eatComment(start int, text []rune) (ok bool) {

	ok = true
	for i := start; i < len(text); i++ {

		r := text[i]
		if unicode.IsSpace(r) {
			continue
		}

		if !CommentMark[r] {
			return false
		}
		return
	}
	return
}
