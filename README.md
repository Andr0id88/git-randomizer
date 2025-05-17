# gitr ― the **Git Randomizer**

> “Why say ‘Fix typo’ when you can say
> ‘Ye olde letter-polish hath been administered, verily!’”
> - *Shakespeare, probably*

> “Why say ‘Update README’ when you can say
> ‘Many people are saying this documentation is the best-tremendous words, the best words!’”
> - *Trump, definitely*

gitr rewrites your **commit messages** and **branch names** in the voice of whoever (or whatever) you fancy:
Yoda, Deadpool, Jim Lahey, Gandalf, Doge, your local weatherman… you name it.
Perfect for side projects, internal tooling, or simply annoying your teammates during code review.

⚠️ **SERIOUS WARNING**
This tool is **satire**. If you push these messages to a production repo your CI may pass,
but HR might not. You’ve been warned.

---

## Features

| Area | What it does |
|------|--------------|
| **Commit messages** | Takes your boring message and runs it through Google Gemini, rewriting it in the chosen persona, mood, and length. Interactive “Yes / generate again / use original / cancel” loop. |
| **Branch names** | Enter a base idea – gitr returns a short/medium kebab-case slug in character (and even `git checkout -b` for you if you approve). |
| **Random everything** | `--random`, `--group cartoons`, or config defaults like `default_mood: random` re-roll persona/mood every generation. |
| **Groups** | Built-in sets: `trailer_park_boys`, `cartoons`, `politicians`, `celebrities`, `literary`, `misc` (plus anything you add). |
| **Tagline** | After a successful commit, adds a one-liner in a separate persona (“Yoda says: Committed, your code is”). |
| **pass secret** | Stores `GEMINI_API_KEY` in `pass` so you don’t leak it into your shell history. |
| **Config + autosave** | First run drops a commented `~/.config/git-randomizer/git-randomizer.yaml`. Use `--save` to write new defaults from CLI flags. |
| **Safety exits** | Press **Ctrl-C** or select **Cancel** at any prompt → immediate clean exit, no commit/branch created. |

---

## Quick start


> ⚠️ **Note:** A Google Gemini API key is required to use `git-randomizer`, regardless of how it's installed.

---

### 🔐 Step 1: Set up your Gemini API key

```bash
# 1. Get your API key from Google AI Studio:
https://aistudio.google.com/apikey

# 2. Export it as an environment variable:
export GEMINI_API_KEY=your-api-key

# (Optional) Store it in pass:
pass insert gemini_api_key

# Then configure your YAML to use it:
# $HOME/.config/gitrandomizer/gitrandomizer.yaml
# (This is set as default, and uses env if it does not exsist)
# pass_secret: gemini_api_key
```
---

### 💾 Step 2: Install git-randomizer

#### 🔧 Option A - 1-liner (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/Andr0id88/git-randomizer/main/install.sh | bash
```
*Installs the gitr binary to ~/.local/bin. Make sure that directory is in your PATH.*

#### 📦 Option B: Download from the Release Page
```bash

# 1. Visit the latest release and download the binary for your OS:
https://github.com/Andr0id88/git-randomizer/releases/latest

# 2. Make the binary executable:
chmod +x gitr

# 3. Move it to a directory in your PATH (e.g.):
mv gitr ~/.local/bin/
```

#### 🛠️ Option C: Build from Source
```bash
# 1. clone and build
git clone https://Andr0id88/git-randomizer.git
cd git-randomizer
go mod tidy && go build -o ~/.local/bin/gitr .
chmod +x ~/.local/bin/gitr
```

---

#### ✅ Step 3: Verify Installation
```bash
gitr --help
```

---

#### 🧪 Step 4: Try It Out!
```bash
cd ~/code/my-project
gitr commit -s "deadpool"
gitr branch -g cartoons -r
```

---

# CLI cheat-sheet
```bash
gitr commit --help   # full list

# Common flags
-s, --style     persona (string or 'random')
-g, --group     pick random persona from group
-r, --random    fully random persona (ignores --style)
-m, --mood      playful | sarcastic | ... | random
-l, --length    short | medium | long
-y, --yes       skip approval step
-p, --pass-secret path/in/pass
-S, --save      write these flags back to YAML defaults

gitr branch [...]   # same vibe, plus:
                    #   generates slug & checks out branch
-L / -G             # list all styles / groups
```

---

# Configuration file

The first run creates:

```bash
# ~/.config/git-randomizer/git-randomizer.yaml
default_character: random        # or "homer simpson"
default_group: ""                # e.g. "cartoons"
default_mood: playful            # or "random"
default_length: medium
confirm: true                    # ask before commit
pass_secret: gemini_api_key      # path in pass
tagline_enabled: true
tagline_persona: yoda
branch_persona: random
branch_persona_group: ""
```

*Change a value or set it to random to enable randomness.*

---

# Adding your own personas 🎭

1. Edit internal/styles/styles.go
```bash
var Groups = map[string][]string{
    "my_buddies": {
        "steve the intern",
        "maria the barista",
        "grandma in caps lock",
    },
    // existing groups...
}
```

2. Re-build

```bash
go vet ./...
go build -o ~/.local/bin/gitr .
```

3. Try it
```bash
gitr commit -g my_buddies -r
gitr commit --list-styles        # your new names appear
gitr commit --list-groups        # group listed too
```

*(No code outside styles.go cares what names you invent, as long as they’re strings.)*

---

# Why does this exist?

Because git history is for humans (and occasionally raccoons).
A touch of humor:

    breaks the monotony of “Refactor util” commits,

    sparks joy in code review,

    reminds you that programming is supposed to be fun.

If your corporate guidelines demand “conventional commits,” run gitr on a personal fork, chuckle, then cherry-pick a sensible message upstream. Everybody wins.

**Have fun — and may your git log read like the wildest cross-over episode never aired.**
