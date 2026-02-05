package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

func OkResponse(req Request, data any) Response {
	return Response{ID: req.ID, Ok: true, Data: MustJSON(data)}
}

func BadResponse(req Request, err error) Response {
	return Response{ID: req.ID, Ok: false, Error: err.Error()}
}

func WriteJSON(w *bufio.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := w.Write(append(b, '\n')); err != nil {
		return err
	}
	return w.Flush()
}

func MustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func TrimLine(b []byte) []byte {
	// removes trailing \n or \r\n
	n := len(b)
	for n > 0 && (b[n-1] == '\n' || b[n-1] == '\r') {
		n--
	}
	return b[:n]
}

func PrettyPrintResponse(resp Response) {
	if resp.Ok {
		fmt.Printf("ok (id=%s)\n", resp.ID)
		if len(resp.Data) > 0 && string(resp.Data) != "null" {
			var pretty any
			if err := json.Unmarshal(resp.Data, &pretty); err == nil {
				out, _ := json.MarshalIndent(pretty, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Println(string(resp.Data))
			}
		}
		return
	}

	fmt.Printf("error (id=%s): %s\n", resp.ID, resp.Error)
	if len(resp.Data) > 0 && string(resp.Data) != "null" {
		fmt.Println(string(resp.Data))
	}
}

func SplitArgs(s string) []string {
	var out []string
	var cur strings.Builder
	inQuotes := false

	flush := func() {
		if cur.Len() > 0 {
			out = append(out, cur.String())
			cur.Reset()
		}
	}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		switch ch {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				cur.WriteByte(ch)
			} else {
				flush()
			}
		default:
			cur.WriteByte(ch)
		}
	}
	flush()
	return out
}
