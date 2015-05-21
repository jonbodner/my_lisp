package assert
import (
    "testing"
    "runtime"
    "strings"
    "fmt"
)

type Assert testing.T

func (a Assert) True(msg string, test bool) {
    if !test {
        a.Fatalf("%s -- True Test failed: %s", getTestLineString(), msg)
    }
}

func (a Assert) Equals(msg string, val1, val2 interface{}) {
    if val1 != val2 {
        a.Fatalf("%s -- Equals Test failed --  %s: Expected '%s', got '%s'", getTestLineString(), msg, val1, val2)
    }
}

func (a Assert) NotNil(msg string, val interface{}) {
    if val == nil {
        a.Fatalf("%s -- Not Nil Test failed -- %s", getTestLineString(), msg)
    }
}

func (a Assert) Nil(msg string, val interface{}) {
    if val != nil {
        a.Fatalf("%s -- Nil Test failed -- %s", getTestLineString(), msg)
    }
}

func getTestLineString() string {
    _, file, line, ok := runtime.Caller(2) // decorate + log + public function.
    if ok {
        // Truncate file name at last file name separator.
        if index := strings.LastIndex(file, "/"); index >= 0 {
            file = file[index+1:]
        } else if index = strings.LastIndex(file, "\\"); index >= 0 {
            file = file[index+1:]
        }
    } else {
        file = "???"
        line = 1
    }
    return fmt.Sprintf("%s:%d", file, line)
}