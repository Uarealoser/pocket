package xlog

import (
	"bytes"
	"strconv"
	"testing"
)

func TestGenLogID(t *testing.T) {
	t.Log(len(strconv.FormatInt(GenLogID(), 10)))
	t.Log(len(DEFAULT_LOG_ID))
	t.Log(GenLogID())
	t.Log(GenLogID())
	t.Log(GenLogID())
	t.Log(GenLogID())
	t.Log(GenLogID())
	t.Log(GenLogID())

}

func TestGetRuntimeInfo(t *testing.T) {
	expected := map[string]interface{}{
		"function": "micode.be.xiaomi.com/kangaroo/pkg/xlog.TestGetRuntimeInfo",
		// "lineno":   10,
	}
	function, _, _ := getRuntimeInfo(1)
	if function != expected["function"] {
		t.Errorf("expected function=%s, get function=%s", expected["function"], function)
	}
}

func TestLogPanic(t *testing.T) {
	logPanic("log panic: %s", "it is a test")
	t.Log("end")
}

func TestPathExists(t *testing.T) {
	cases := []map[string]bool{
		{"/aa/aa": false},
		{"/home/": true},
		{"./xxx": false},
	}
	for _, c := range cases {
		for path, expected := range c {
			e, err := pathExists(path)
			if err != nil {
				t.Errorf("path:%s Exists err=%v", path, err)
			}
			if e != expected {
				t.Errorf("expected %v, get %v", expected, e)
			}
		}
	}
}
func TestFormat(t *testing.T) {
	buffer := new(bytes.Buffer)
	runtimeFormat(buffer, "test", "testb", 12)

	got := string(buffer.Bytes())
	want := "test:testb:12"
	if got != want {
		t.Errorf("got:%s,want:%s", got, want)
	}

	buffer.Reset()
	tagsFormat(buffer, "aaa", "bbb", "ccc")
	got = string(buffer.Bytes())
	want = "[aaa] [bbb] [ccc] "
	if got != want {
		t.Errorf("got:%s,want:%s", got, want)
	}
}

func TestFieldFormat(t *testing.T) {
	buffer := new(bytes.Buffer)

	fieldsNormalFormat(buffer, "it is a message", map[string]string{"a": "b"})
	got := string(buffer.Bytes())
	t.Log(got)
	buffer = new(bytes.Buffer)
	fieldsJsonFormat(buffer, map[string]string{"a": "b"})
	got = string(buffer.Bytes())
	t.Log(got)
	// want = `[order_id:xxx,price:123.23,message:it is a message,user_id:123] it is a message`
	// if got != want {
	// 	t.Errorf("want:%s,got:%s", want, got)
	// }
}
