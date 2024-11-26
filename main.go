package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"nvim-http/lsp"
	"nvim-http/rpc"
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
		err = WriteResponse(lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				Uri: request.Params.TextDocument.Uri,
				Diagnostics: []lsp.Diagnostic{
					{
						Range: lsp.Range{
							Start: lsp.Position{
								Line:      0,
								Character: 0,
							},
							End: lsp.Position{
								Line:      0,
								Character: 5,
							},
						},
						Severity: 1,
						Source:   "me stuff",
						Message:  "you made a change",
					},
					{
						Range: lsp.Range{
							Start: lsp.Position{
								Line:      1,
								Character: 0,
							},
							End: lsp.Position{
								Line:      1,
								Character: 5,
							},
						},
						Severity: 1,
						Source:   "me stuff 2",
						Message:  "you made a change",
					},
				},
			},
		})
		if err != nil {
			logger.Fatalf("error, when WriteResponse() for textDocument/didChange method for handleMessage(). Error: %v", err)
		}
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
			Result: []lsp.CompletionItem{
				{
					Label:         "NVIM (btw)",
					Detail:        "Best editor",
					Documentation: "the only editor I use",
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
