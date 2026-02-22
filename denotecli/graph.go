// graph.go
package main

import (
	"os"
	"strings"
)

// GraphLink represents a link between two notes.
type GraphLink struct {
	SourceID    string `json:"source_id"`
	SourceTitle string `json:"source_title"`
	TargetID    string `json:"target_id"`
}

// GraphResult holds graph traversal results for a note.
type GraphResult struct {
	DenoteFile
	Outgoing []GraphLink `json:"outgoing"`
	Incoming []GraphLink `json:"incoming"`
}

// BuildGraph returns outgoing and incoming links for a given note ID.
func BuildGraph(files []DenoteFile, id string, depth int) (GraphResult, error) {
	// Build ID → file index
	index := make(map[string]*DenoteFile, len(files))
	for i := range files {
		index[files[i].ID] = &files[i]
	}

	target, ok := index[id]
	if !ok {
		return GraphResult{}, errNotFound(id)
	}

	// Outgoing: links from this note
	data, err := os.ReadFile(target.Path)
	if err != nil {
		return GraphResult{}, err
	}
	outIDs := ExtractLinks(string(data))
	var outgoing []GraphLink
	for _, oid := range outIDs {
		link := GraphLink{SourceID: id, SourceTitle: target.Title, TargetID: oid}
		outgoing = append(outgoing, link)
	}

	// Incoming: other notes that link to this ID
	var incoming []GraphLink
	needle := "denote:" + id
	for i := range files {
		if files[i].ID == id {
			continue
		}
		d, err := os.ReadFile(files[i].Path)
		if err != nil {
			continue
		}
		if strings.Contains(string(d), needle) {
			incoming = append(incoming, GraphLink{
				SourceID:    files[i].ID,
				SourceTitle: files[i].Title,
				TargetID:    id,
			})
		}
	}

	// Enrich frontmatter for target
	fm := ParseFrontmatter(string(data))
	if fm.Title != "" {
		target.Title = fm.Title
	}
	if len(fm.Filetags) > 0 {
		target.Tags = fm.Filetags
	}

	return GraphResult{
		DenoteFile: *target,
		Outgoing:   outgoing,
		Incoming:   incoming,
	}, nil
}

func errNotFound(id string) error {
	return &notFoundError{id}
}

type notFoundError struct{ id string }

func (e *notFoundError) Error() string { return "not found: " + e.id }
