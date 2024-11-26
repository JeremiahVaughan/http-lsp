package lsp


type DidChangeTextDocumentNotification struct {
    Notification
    Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
    TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
        ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

/**
 * An event describing a change to a text document. If only a text is provided
 * it is considered to be the full content of the document.
 */
type TextDocumentContentChangeEvent struct {
    Text string `json:"text"`
};

