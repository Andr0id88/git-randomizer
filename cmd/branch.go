package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"git-randomizer/internal/gemini"
	"git-randomizer/internal/styles"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/* ------------------------------ FLAGS ------------------------------ */

var (
	brStyle      string
	brRandom     bool
	brGroup      string
	brMood       string
	brLength     string
	brPass       string
	brListGroups bool
	brSave       bool
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Generate a punny branch name in character and optionally check it out",
	RunE:  runBranch,
}

func init() {
	branchCmd.Flags().StringVarP(&brStyle, "style", "s", "", "persona (or 'random')")
	branchCmd.Flags().BoolVarP(&brRandom, "random", "r", false, "fully random persona")
	branchCmd.Flags().StringVarP(&brGroup, "group", "g", "", "random persona from this group")
	branchCmd.Flags().StringVarP(&brMood, "mood", "m", "", "mood or 'random'")
	branchCmd.Flags().StringVarP(&brLength, "length", "l", "short", "short | medium")
	branchCmd.Flags().StringVarP(&brPass, "pass-secret", "p", "", "pass secret for GEMINI_API_KEY")
	branchCmd.Flags().BoolVarP(&brListGroups, "list-groups", "G", false, "list persona groups & exit")
	branchCmd.Flags().BoolVarP(&brSave, "save", "S", false, "save persona/group defaults")
}

/* ---------------------------- COMMAND ----------------------------- */

func runBranch(cmd *cobra.Command, _ []string) error {
	if brListGroups {
		fmt.Println("Available groups:")
		for _, g := range styles.GroupNames() {
			fmt.Printf("  • %s\n", g)
		}
		return nil
	}

	if _, err := os.Stat(".git"); err != nil {
		return errors.New("❌ not inside a git repository")
	}

	apiKey, err := getAPIKey(brPass)
	if err != nil {
		return err
	}

	base, err := promptBaseName()
	if err != nil {
		return err
	}

	persona := pickPersona()
	mood := pickBranchMoodOnce()
	lengthRule := "short"
	if strings.ToLower(brLength) == "medium" {
		lengthRule = "medium"
	}

	for {
		slug, err := generateSlug(apiKey, base, persona, mood, lengthRule)
		if err != nil {
			return err
		}
		fmt.Printf("\n🌿 Suggested branch (%s, %s): %s\n\n", persona, mood, slug)

		/* -------- Confirmation prompt -------- */
		confirm := promptui.Prompt{
			Label:     "✅ Use this branch name?",
			IsConfirm: true,
			Default:   "Y",
		}
		ans, cerr := confirm.Run()
		switch cerr {
		case promptui.ErrInterrupt, promptui.ErrEOF:
			fmt.Println("\n🚫 Aborted.")
			return nil
		case promptui.ErrAbort:
			// user typed "n" → treat as No, fall through to menu
		case nil:
			if ans == "" || strings.ToLower(ans) == "y" {
				if err := checkoutBranch(slug); err != nil {
					return err
				}
				fmt.Println("✅ Switched to new branch!")
				return nil
			}
		default:
			return cerr
		}

		/* -------- Secondary menu -------- */
		menu := promptui.Select{
			Label:        "❓ What next?",
			Items:        []string{"Generate another", "Use my original text", "Cancel"},
			HideSelected: true,
		}
		_, act, serr := menu.Run()
		if serr == promptui.ErrInterrupt || serr == promptui.ErrEOF {
			fmt.Println("\n🚫 Aborted.")
			return nil
		} else if serr != nil {
			return serr
		}

		switch act {
		case "Generate another":
			if personaIsRandom() {
				persona = pickPersona()
			}
			if branchMoodIsRandom() {
				mood = styles.RandomMood()
			}
			continue
		case "Use my original text":
			slug = slugify(base)
			if err := checkoutBranch(slug); err != nil {
				return err
			}
			fmt.Println("✅ Switched to new branch!")
			return nil
		default: /* Cancel */
			fmt.Println("🚫 Aborted.")
			return nil
		}
	}
}

/* ----------------------- LOGIC HELPERS ---------------------- */

func personaIsRandom() bool {
	if brRandom || strings.ToLower(brStyle) == "random" {
		return true
	}
	if brGroup != "" && brStyle == "" {
		return true
	}
	if strings.ToLower(viper.GetString("branch_persona")) == "random" {
		return true
	}
	if grp := viper.GetString("branch_persona_group"); grp != "" && brStyle == "" && brGroup == "" {
		return true
	}
	return false
}

func branchMoodIsRandom() bool {
	if strings.ToLower(brMood) == "random" {
		return true
	}
	return strings.ToLower(viper.GetString("default_mood")) == "random"
}

func pickPersona() string {
	if !personaIsRandom() {
		if brStyle != "" && strings.ToLower(brStyle) != "random" {
			return brStyle
		}
		if p := viper.GetString("branch_persona"); p != "" && strings.ToLower(p) != "random" {
			return p
		}
	}
	if brGroup != "" {
		if s, ok := styles.RandomFromGroup(brGroup); ok {
			return s
		}
	}
	if brRandom || strings.ToLower(brStyle) == "random" {
		return styles.Random()
	}
	if grp := viper.GetString("branch_persona_group"); grp != "" {
		if s, ok := styles.RandomFromGroup(grp); ok {
			return s
		}
	}
	return styles.Random()
}

func pickBranchMoodOnce() string {
	if !branchMoodIsRandom() {
		if brMood != "" && strings.ToLower(brMood) != "random" {
			return brMood
		}
		if def := viper.GetString("default_mood"); def != "" && strings.ToLower(def) != "random" {
			return def
		}
	}
	return styles.RandomMood()
}

/* ----------------------- I/O & GIT ------------------------- */

func promptBaseName() (string, error) {
	fmt.Print("📝 Base branch description: ")
	in := bufio.NewReader(os.Stdin)
	txt, err := in.ReadString('\n')
	return strings.TrimSpace(txt), err
}

func generateSlug(key, base, persona, mood, length string) (string, error) {
	prompt := fmt.Sprintf(
		`Rewrite the text below as a very short git branch slug in the style of %s with a %s vibe. Use kebab-case. Keep it %s (max 40 chars). Respond with the slug only.
Text:
"""%s"""`, persona, mood, length, base)

	out, err := gemini.Generate(key, persona, mood, length, prompt)
	if err != nil {
		return "", err
	}
	return slugify(out), nil
}

func slugify(in string) string {
	s := strings.ToLower(in)
	s = regexp.MustCompile(`[^a-z0-9\s\-]`).ReplaceAllString(s, " ")
	s = strings.Join(strings.Fields(strings.ReplaceAll(s, "_", " ")), "-")
	if len(s) > 40 {
		s = s[:40]
	}
	return strings.Trim(s, "-")
}

func checkoutBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

