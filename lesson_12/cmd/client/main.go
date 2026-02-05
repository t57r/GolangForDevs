package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"lesson12/internal/model"
)

var reqSeq uint64

func nextID() string {
	n := atomic.AddUint64(&reqSeq, 1)
	return strconv.FormatUint(n, 10)
}

func main() {
	var addr = model.Addr
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("connect error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", addr)
	fmt.Println("Type 'help' for commands. Type 'exit' to quit.")

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("stdin error: %v\n", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			return
		}
		if line == "help" {
			printHelp()
			continue
		}

		req, err := parseCommand(line)
		if err != nil {
			fmt.Printf("command error: %v\n", err)
			continue
		}

		// send request
		if err := model.WriteJSON(w, req); err != nil {
			fmt.Printf("send error: %v\n", err)
			return
		}

		// wait response
		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		respLine, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Printf("read response error: %v\n", err)
			return
		}
		_ = conn.SetReadDeadline(time.Time{})

		var resp model.Response
		if err := json.Unmarshal(model.TrimLine(respLine), &resp); err != nil {
			fmt.Printf("bad response json: %v\nraw=%s\n", err, string(respLine))
			continue
		}

		model.PrettyPrintResponse(resp)
	}
}

func printHelp() {
	fmt.Println(`
Commands:

# collections
create-collection <name> <primaryKey>
get-collection <name>
delete-collection <name>

# documents
put <collection> <id> <name>
get <collection> <id>
delete <collection> <id>
list <collection>

# indexes
create-index <collection> <fieldName>
delete-index <collection> <fieldName>

# query (range query on string field)
# min/max are optional, use "-" to mean nil
query <collection> <fieldName> <asc|desc> <min|- > <max|- >

Examples:
create-collection users ID
put users 1 Alice
put users 2 Bob
create-index users Name
query users Name asc A Z
get users 1
list users
delete users 2
delete-collection users
`)
}

func parseCommand(line string) (*model.Request, error) {
	parts := model.SplitArgs(line)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := parts[0]
	id := nextID()

	switch cmd {

	// --- collection ops ---
	case "create-collection":
		if len(parts) != 3 {
			return nil, fmt.Errorf("usage: create-collection <name> <primaryKey>")
		}
		name := parts[1]
		pk := parts[2]
		payload := model.MustJSON(map[string]any{"primaryKey": pk})
		return &model.Request{ID: id, Op: model.RequestOperationCreateCollection, Collection: name, Payload: payload}, nil

	case "get-collection":
		if len(parts) != 2 {
			return nil, fmt.Errorf("usage: get-collection <name>")
		}
		return &model.Request{ID: id, Op: model.RequestOperationGetCollection, Collection: parts[1]}, nil

	case "delete-collection":
		if len(parts) != 2 {
			return nil, fmt.Errorf("usage: delete-collection <name>")
		}
		return &model.Request{ID: id, Op: model.RequestOperationDeleteCollection, Collection: parts[1]}, nil

	// --- document ops ---
	case "put":
		// put <collection> <id> <name>
		if len(parts) < 4 {
			return nil, fmt.Errorf("usage: put <collection> <id> <name>")
		}
		collection := parts[1]
		docID := parts[2]
		name := strings.Join(parts[3:], " ")

		doc := map[string]any{
			"Fields": map[string]any{
				"ID":   map[string]any{"Type": "string", "Value": docID},
				"Name": map[string]any{"Type": "string", "Value": name},
			},
		}
		payload := model.MustJSON(map[string]any{"doc": doc})
		return &model.Request{ID: id, Op: model.RequestOperationPutDocument, Collection: collection, Payload: payload}, nil

	case "get":
		// get <collection> <id>
		if len(parts) != 3 {
			return nil, fmt.Errorf("usage: get <collection> <id>")
		}
		payload := model.MustJSON(map[string]any{"key": parts[2]})
		return &model.Request{ID: id, Op: model.RequestOperationGetDocument, Collection: parts[1], Payload: payload}, nil

	case "delete":
		// delete <collection> <id>
		if len(parts) != 3 {
			return nil, fmt.Errorf("usage: delete <collection> <id>")
		}
		payload := model.MustJSON(map[string]any{"key": parts[2]})
		return &model.Request{ID: id, Op: model.RequestOperationDeleteDocument, Collection: parts[1], Payload: payload}, nil

	case "list":
		if len(parts) != 2 {
			return nil, fmt.Errorf("usage: list <collection>")
		}
		return &model.Request{ID: id, Op: model.RequestOperationListDocuments, Collection: parts[1]}, nil

	// --- index ops ---
	case "create-index":
		if len(parts) != 3 {
			return nil, fmt.Errorf("usage: create-index <collection> <fieldName>")
		}
		payload := model.MustJSON(map[string]any{"fieldName": parts[2]})
		return &model.Request{ID: id, Op: model.RequestOperationCreateIndex, Collection: parts[1], Payload: payload}, nil

	case "delete-index":
		if len(parts) != 3 {
			return nil, fmt.Errorf("usage: delete-index <collection> <fieldName>")
		}
		payload := model.MustJSON(map[string]any{"fieldName": parts[2]})
		return &model.Request{ID: id, Op: model.RequestOperationDeleteIndex, Collection: parts[1], Payload: payload}, nil

	case "query":
		// query <collection> <fieldName> <asc|desc> <min|- > <max|- >
		if len(parts) != 6 {
			return nil, fmt.Errorf("usage: query <collection> <fieldName> <asc|desc> <min|- > <max|- >")
		}
		collection := parts[1]
		field := parts[2]
		order := strings.ToLower(parts[3])

		desc := false
		switch order {
		case "asc":
			desc = false
		case "desc":
			desc = true
		default:
			return nil, fmt.Errorf("order must be asc or desc")
		}

		var minPtr *string
		var maxPtr *string

		if parts[4] != "-" {
			v := parts[4]
			minPtr = &v
		}
		if parts[5] != "-" {
			v := parts[5]
			maxPtr = &v
		}

		payload := model.MustJSON(map[string]any{
			"fieldName": field,
			"desc":      desc,
			"minValue":  minPtr,
			"maxValue":  maxPtr,
		})
		return &model.Request{ID: id, Op: model.RequestOperationQuery, Collection: collection, Payload: payload}, nil

	default:
		return nil, fmt.Errorf("unknown command %q (type 'help')", cmd)
	}
}
