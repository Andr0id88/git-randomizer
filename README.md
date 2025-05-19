# gitr - the **Git Randomizer**

Why settle for:

> `fix: typo`

When you could inscribe your legacy into the `git log` like:

---

> “Ye olde letter-polish hath been administered, verily!”
>
> - *Shakespeare, probably*

---

> ‘Many people are saying this documentation is the best - tremendous words, the best words!’”
>
> - *Trump, definitely*

---

Your `git log` should feel like a **cursed scroll** -
etched in arcane symbols, whispered about in hushed tones by future interns.

This isn’t just version control.
It’s a **ritual**, a **rebellion**, a **beautiful mistake** waiting to be discovered.

**Write weird. Commit louder.
Let the chaos compile.**

---

### 🌀 What is this?

`gitr` rewrites your **commit messages** and **branch names** in the voice of whoever (or *whatever*) you fancy:
Yoda, Michael Jackson, Deadpool, Jim Lahey, Gandalf, Doge, your local weatherman… you name it.

Perfect for side projects, internal tooling, or simply annoying your teammates during code review.

⚠️ **SERIOUS WARNING**
This tool is *satire*. If you push these messages to a production repo, your CI may pass -
but HR might not. You’ve been warned.

---

## ✨ Features

| Area | What it does |
|------|--------------|
| **Commit messages** | Takes your boring message and runs it through Google Gemini, rewriting it in the chosen persona, mood, and length. Interactive “Yes / generate again / use original / cancel” loop. |
| **Branch names** | Enter a base idea – gitr returns a short/medium kebab-case slug in character (and even `git checkout -b` for you if you approve). |
| **Random everything** | `--random`, `--group cartoons`, or config defaults like `default_mood: random` re-roll persona/mood every generation. |
| **Groups** | Built-in sets: `supervillains`, `cartoons`, `politicians`, `celebrities`, `conspiracy_theorists`, `misc` (plus many more and anything you add). |
| **Tagline** | After a successful commit, adds a one-liner in a separate persona (“Yoda says: Committed, your code is”). |
| **pass secret** | Stores `GEMINI_API_KEY` in `pass` so you don’t leak it into your shell history. |
| **Config + autosave** | First run drops a commented `~/.config/git-randomizer/git-randomizer.yaml`. Use `--save` to write new defaults from CLI flags. |
| **Safety exits** | Press **Ctrl-C** or select **Cancel** at any prompt → immediate clean exit, no commit/branch created. |

---

## 🚀 Quick start


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
gitr commit -s "Donald Trump"
gitr branch -g cartoons -r
```

---

## 💻 CLI cheat-sheet
```bash
gitr commit --help   # full list

# Common flags
-s, --style       persona (string or 'random')
-g, --group       pick random persona from group
-r, --random      fully random persona (ignores --style)
-m, --mood        playful | sarcastic | ... | random
-l, --length      short | medium | long
-y, --yes         skip approval step
-p, --pass-secret path/in/pass
-S, --save        write these flags back to YAML defaults
-L / -G           list all styles / groups

gitr branch [...]   # same vibe, plus: generates slug & checks out branch
```

---

## 📝 Configuration file

The first run creates:

```bash
# ------------------------------------------------------------
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
branch_persona_group: ""        # e.g. "supervillains"
```

*Change a value or set it to random to enable randomness.*

---

## 🎭 Adding your own personas

1. Edit internal/styles/styles.go
```bash
var Groups = map[string][]string{
    "my_buddies": {
        "grandma in caps lock",
        "steve the intern",
        "maria the barista",
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

## 🤔 Why does this exist?

Because the world has enough boring commit messages -
and not nearly enough that say things like:

---

> “Screw you guys, I’m rebasing at home.”
>
> - *Eric Cartman, 2 seconds before force-pushing to `main`*

---

> “After three hours of debugging, I realized the bug was me all along.”
>
> - Sir Galahad, softly weeping into main.py

---

### So really... why?

- Because `git log` should read like the **diary of a caffeinated raccoon**.
- Because you’ve always wanted to `git blame` a bug on **Alex Jones** screaming about interdimensional bugs introduced by globalist devs.
- Because every repo deserves a branch named `feature/screw-you-guys-im-going-home`.
- Because you once typed `fix thing` and felt a part of you die inside.
- Most importantly, because programming is supposed to be **fun.**

---

Yes, you *could* follow conventional commit standards…
But wouldn’t you rather write something with a little more ✨ chaos?

If your corporate guidelines demand “conventional commits,”
go full goblin in a personal fork, laugh, then cherry-pick a respectable message upstream.
**Professional on the outside, feral on the inside.**

---

Use it. Abuse it. Confuse your team.
Fill your repo with mystery, mayhem, and markdown-stained madness.

Because someday, archaeologists will unearth your git history and whisper,
> “What the hell were they building?”

And that, my friend, is **legacy.**

---

🤝 Contributing (Optional)

Contributions are welcome! If you have ideas for improvements or find a bug, please feel free to:

1. Fork the repository.
2. Clone your fork to your local machine:
```bash
git clone https://github.com/andr0id88/git-randomizer.git
```
3. Create a new feature branch:
```bash
git checkout -b feature/AmazingFeature
```
4. Make your changes.
5. Commit and push your changes:
```bash
git commit -m 'Add some AmazingFeature'
git push origin feature/AmazingFeature
```
6. Open a Pull Request.

---

📄 License
This project is licensed under the [MIT License](https://github.com/Andr0id88/git-randomizer/blob/main/LICENSE)
