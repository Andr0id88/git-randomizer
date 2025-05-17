package config

import (
	"os"
	"path/filepath"
)

/*
CreateDefault writes an example-rich YAML file if none exists.

We can’t make Viper emit comments, so we craft the file manually here on
first run.  Users can delete values or set them to “random” to enable
randomisation.
*/
func CreateDefault(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	sample := `# ------------------------------------------------------------
# git-randomizer configuration -- every CLI flag has a YAML twin.
# Anything here is overridden by command-line flags at runtime.
# Set a value to 'random' to enable randomisation.
# ------------------------------------------------------------

# --- Commit defaults ----------------------------------------
default_character: random   # persona, e.g. "yoda" or "donald trump"
default_group: ""           # e.g. "cartoons" – random within group
default_mood: playful       # 'playful', 'sarcastic', or 'random'
default_length: medium      # short | medium | long
confirm: true               # true = ask before committing

# --- API key storage ----------------------------------------
pass_secret: "gemini_api_key"   # path in 'pass' – overrides GEMINI_API_KEY

# --- Success tagline ----------------------------------------
tagline_enabled: true
tagline_persona: yoda           # persona for the one-liner after commit

# --- Branch-name generator ----------------------------------
branch_persona: random          # fixed persona OR 'random'
branch_persona_group: ""        # e.g. "trailer_park_boys"
`

	return os.WriteFile(path, []byte(sample), 0o644)
}

