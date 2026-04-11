// SPDX-License-Identifier: AGPL-3.0-or-later

package vectorstore

import (
	"strings"
	"unicode"
)

// TextChunk represents a chunk of text produced by the tokenizer.
type TextChunk struct {
	Text       string
	TokenCount int
}

// tokenize splits text into tokens on whitespace and punctuation boundaries.
// Each contiguous run of letters/digits is a token. This is a simple
// approximation — it doesn't match any specific model's BPE tokenizer, but is
// good enough for determining chunk boundaries.
func tokenize(text string) []string {
	var tokens []string
	var cur strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cur.WriteRune(r)
		} else {
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			// Keep meaningful punctuation as separate tokens.
			if !unicode.IsSpace(r) {
				tokens = append(tokens, string(r))
			}
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

// ChunkText splits text into overlapping chunks of at most maxTokens tokens.
// overlap specifies how many tokens consecutive chunks share. If the entire
// text fits in one chunk, a single TextChunk is returned.
func ChunkText(text string, maxTokens, overlap int) []TextChunk {
	if maxTokens <= 0 {
		maxTokens = 2048
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= maxTokens {
		overlap = maxTokens / 10
	}

	tokens := tokenize(text)
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) <= maxTokens {
		return []TextChunk{{Text: text, TokenCount: len(tokens)}}
	}

	step := maxTokens - overlap
	if step <= 0 {
		step = 1
	}

	var chunks []TextChunk
	for start := 0; start < len(tokens); start += step {
		end := min(start+maxTokens, len(tokens))
		chunkTokens := tokens[start:end]
		chunkText := strings.Join(chunkTokens, " ")
		chunks = append(chunks, TextChunk{
			Text:       chunkText,
			TokenCount: len(chunkTokens),
		})
		if end == len(tokens) {
			break
		}
	}
	return chunks
}
