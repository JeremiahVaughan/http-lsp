package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/JeremiahVaughan/http-lsp/lsp"
	"github.com/JeremiahVaughan/http-lsp/rpc"
	"log"
	"os"
)

func main() {
	logger := getLogger("/tmp/nvim-http-log.txt")
	logger.Println("starting")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, bytes, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Fatalf("error, when rpc.DecodeMessage() for main(). Error: %v", err)
		}
		err = handleMessage(logger, method, bytes)
		if err != nil {
			logger.Fatalf("error, when handleMessage() for main(). Error: %v", err)
		}
	}
}

func handleMessage(logger *log.Logger, method string, bytes []byte) error {
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal initalize request for handleMessage(). Error: %v", err)
		}
		logger.Printf(
			"connected to %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(request.ID)
		err = WriteResponse(msg)
		if err != nil {
			logger.Fatalf("error, when WriteResponse() for initalize method for handleMessage(). Error: %v", err)
		}
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/didOpen request for handleMessage(). Error: %v", err)
		}
		logger.Printf(
			"Opened: %s",
			request.Params.TextDocument.Uri,
		)
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/didChange request for handleMessage(). Error: %v", err)
		}
		logger.Printf(
			"Changed: %s",
			request.Params.TextDocument.Uri,
		)
		// err = WriteResponse(lsp.PublishDiagnosticsNotification{
		// 	Notification: lsp.Notification{
		// 		RPC:    "2.0",
		// 		Method: "textDocument/publishDiagnostics",
		// 	},
		// 	Params: lsp.PublishDiagnosticsParams{
		// 		Uri: request.Params.TextDocument.Uri,
		// 		Diagnostics: []lsp.Diagnostic{
		// 			{
		// 				Range: lsp.Range{
		// 					Start: lsp.Position{
		// 						Line:      0,
		// 						Character: 0,
		// 					},
		// 					End: lsp.Position{
		// 						Line:      0,
		// 						Character: 5,
		// 					},
		// 				},
		// 				Severity: 1,
		// 				Source:   "me stuff",
		// 				Message:  "you made a change",
		// 			},
		// 			{
		// 				Range: lsp.Range{
		// 					Start: lsp.Position{
		// 						Line:      1,
		// 						Character: 0,
		// 					},
		// 					End: lsp.Position{
		// 						Line:      1,
		// 						Character: 5,
		// 					},
		// 				},
		// 				Severity: 1,
		// 				Source:   "me stuff 2",
		// 				Message:  "you made a change",
		// 			},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	logger.Fatalf("error, when WriteResponse() for textDocument/didChange method for handleMessage(). Error: %v", err)
		// }
	case "textDocument/hover":
		var request lsp.HoverRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/hover request for handleMessage(). Error: %v", err)
		}
		resp := lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				Id:  &request.ID,
			},
			Result: lsp.HoverResult{
				Contents: "hello from lsp",
			},
		}
		err = WriteResponse(resp)
		if err != nil {
			return fmt.Errorf("error, when WriteResponse() for textDocument/hover request for handleMessage(). Error: %v", err)
		}
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/definition request for handleMessage(). Error: %v", err)
		}
		resp := lsp.DefinitionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				Id:  &request.ID,
			},
			Result: lsp.Location{
				Uri: request.Params.TextDocument.Uri,
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      request.Params.Position.Line - 1,
						Character: 0,
					},
					End: lsp.Position{
						Line:      request.Params.Position.Line - 1,
						Character: 0,
					},
				},
			},
		}
		err = WriteResponse(resp)
		if err != nil {
			return fmt.Errorf("error, when WriteResponse() for textDocument/definition request for handleMessage(). Error: %v", err)
		}
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/codeAction request for handleMessage(). Error: %v", err)
		}
		resp := lsp.CodeActionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				Id:  &request.ID,
			},
			Result: lsp.CodeActionResult{
				Title: "TODO implement",
			},
		}
		err = WriteResponse(resp)
		if err != nil {
			return fmt.Errorf("error, when WriteResponse() for textDocument/codeAction request for handleMessage(). Error: %v", err)
		}
	case "textDocument/completion":
		var request lsp.CompletionRequest
		err := json.Unmarshal(bytes, &request)
		if err != nil {
			return fmt.Errorf("error, when attempting to unmarshal textDocument/completion request for handleMessage(). Error: %v", err)
		}
		resp := lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				Id:  &request.ID,
			},
			// todo make these values contextually aware
			Result: []lsp.CompletionItem{
				{
					Label:         "http://",
					Detail:        "http",
					Documentation: "unencrypted",
				},
				{
					Label:         "https://",
					Detail:        "https",
					Documentation: "encrypted",
				},
				{
					Label:         "POST ",
					Detail:        "post method",
					Documentation: "post method",
				},
				{
					Label:         "PUT ",
					Detail:        "put method",
					Documentation: "put method",
				},
				{
					Label:         "PATCH ",
					Detail:        "patch method",
					Documentation: "patch method",
				},
				{
					Label:         "GET ",
					Detail:        "get method",
					Documentation: "get method",
				},
				{
					Label:         "DELETE ",
					Detail:        "delete method",
					Documentation: "delete method",
				},
				{
					Label:         "User-Agent: ",
					Detail:        "http header user agent",
					Documentation: "http header user agent",
				},
				{
					Label:         "Accept-Language: ",
					Detail:        "http header accept language",
					Documentation: "http header accept language",
				},
				{
					Label:         "Accept-Encoding: ",
					Detail:        "http header accept encoding",
					Documentation: "http header accept encoding",
				},
				{
					Label:         "br",
					Detail:        "accept-encoding value br",
					Documentation: "accept-encoding value br",
				},
				{
					Label:         "gzip",
					Detail:        "accept-encoding value gzip",
					Documentation: "accept-encoding value gzip",
				},
				{
					Label:         "deflate",
					Detail:        "accept-encoding value deflate",
					Documentation: "accept-encoding value deflate",
				},
				{
					Label:         "Accept: ",
					Detail:        "http header accept",
					Documentation: "http header accept",
				},
				{
					Label:         "Content-Type: ",
					Detail:        "http header content type",
					Documentation: "http header content type",
				},
				{
					Label:         "application/x-www-form-urlencoded",
					Detail:        "content type value",
					Documentation: "content type value",
				},
				{
					Label:         "application/json",
					Detail:        "accept or content type value",
					Documentation: "accept or content type value",
				},
				{
					Label:         "text/html",
					Detail:        "accept value",
					Documentation: "accept value",
				},
				{
					Label:         "Referer: ",
					Detail:        "http header referer",
					Documentation: "http header referer",
				},
				{
					Label:         "Connection: ",
					Detail:        "http header connection",
					Documentation: "http header connection",
				},
				{
					Label:         "Cache-Control: ",
					Detail:        "http header cache-control",
					Documentation: "http header cache-control",
				},
				{
					Label:         "max-age=",
					Detail:        "cache-control value",
					Documentation: "unit is seconds",
				},
				{
					Label:         "Host: ",
					Detail:        "http header host",
					Documentation: "http header host",
				},
				{
					Label:         "HTTP/1.1",
					Detail:        "protocol version 1",
					Documentation: "protocol version version 1",
				},
				{
					Label:         "HTTP/2",
					Detail:        "protocol version 2",
					Documentation: "protocol version version 2",
				},
			},
		}
		err = WriteResponse(resp)
		if err != nil {
			return fmt.Errorf("error, when WriteResponse() for textDocument/codeAction request for handleMessage(). Error: %v", err)
		}
	}
	logger.Printf("received message with method: %s", method)
	return nil
}

func WriteResponse(msg any) error {
	reply, err := rpc.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("error, when rpc.EncodeMessage() for handleMessage(). Error: %v", err)
	}
	writer := os.Stdout
	writer.Write(reply)
	return nil
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didnt give me a good file")
	}

	return log.New(logfile, "[nvim-http]", log.Ldate|log.Ltime|log.Lshortfile)
}
