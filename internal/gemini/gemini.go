package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const model = "gemini-2.0-flash"

type apiReq struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type apiResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Generate calls Gemini and returns a rewritten commit message.
func Generate(apiKey, style, mood, length, commit string) (string, error) {
	prompt := buildPrompt(commit, style, mood, length)

	body := apiReq{Contents: []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{{
		Parts: []struct {
			Text string `json:"text"`
		}{{
			Text: prompt,
		}},
	}}}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model, apiKey)

	httpResp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(httpResp.Body)
		return "", fmt.Errorf("Gemini API error: %s", string(raw))
	}

	var r apiResp
	if err := json.NewDecoder(httpResp.Body).Decode(&r); err != nil {
		return "", err
	}

	if len(r.Candidates) == 0 || len(r.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}
	return strings.TrimSpace(r.Candidates[0].Content.Parts[0].Text), nil
}

func buildPrompt(msg, style, mood, length string) string {
	lengthRule := map[string]string{
		"short":  "Keep it to MAX 8–12 words.",
		"medium": "Aim for one punchy line (≤ 20 words).",
		"long":   "You may use up to ~40 words (two concise lines).",
	}
	rule := lengthRule[length]

	translate := ""
	if strings.Contains(strings.ToLower(style), "ivar aasen") {
		translate = " Translate the commit message into contemporary Nynorsk (New Norwegian) before applying the persona."
	}

	return fmt.Sprintf(
		`Rewrite the following git commit message in the style of %s with a %s mood.%s %s
Respond ONLY with the final rewritten git commit message itself – no pre-amble, no bullet points, no code fences.

Commit message:
"""%s"""`,
		style, mood, translate, rule, msg)
}

