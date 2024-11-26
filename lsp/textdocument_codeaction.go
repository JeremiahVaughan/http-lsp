package lsp


type CodeActionRequest struct {
    Request
    Params CodeActionParams `json:"params"`
}

type CodeActionParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Range Range `json:"range"`
    Context CodeActionContext `json:"context"`
}

type CodeActionContext struct {
}

type CodeActionResponse struct {
    Response
    Result CodeActionResult `json:"result"`
}

type CodeActionResult struct {
    Title string `json:"title"`
    Edit *WorkspaceEdit `json:"edit,omitempty"`
    Command *Command `json:"command,omitempty"`
}

type Command struct {
    Title string `json:"title"`
    Command string `json:"command"`
    Arguments []any `json:"arguments,omitempty"`
}
