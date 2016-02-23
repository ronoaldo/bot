package bot

// History retains in-memory records of navigation entries.
type History struct {
	history []string
}

// Entries return the recent URLs visited.
func (h *History) Entries() []string {
	return h.history
}

// Current returs the most recent visited URL.
// This method is not concurrent safe.
func (h *History) Current() string {
	if len(h.history) == 0 {
		return ""
	}
	return h.history[len(h.history)-1]
}

// Add appends a new entry to the history.
// This method is not concurrent safe.
func (h *History) Add(entry string) {
	h.history = append(h.history, entry)
}
