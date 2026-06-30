package authorization

type Request struct {
    SubjectID string   `json:"subjectId"`
    Roles     []string `json:"roles"`
    Action    string   `json:"action"`
    Resource  string   `json:"resource"`
}

type Response struct {
    Allowed           bool     `json:"allowed"`
    PoliciesEvaluated []string `json:"policiesEvaluated,omitempty"`
    Message           string   `json:"message"`
    Reason            string   `json:"reason,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
