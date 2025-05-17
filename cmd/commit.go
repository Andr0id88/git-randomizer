package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"git-randomizer/internal/gemini"
	"git-randomizer/internal/styles"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/* ---------------------- FLAGS ---------------------- */

var (
	flagStyle      string
	flagRandom     bool
	flagGroup      string
	flagMood       string
	flagLength     string
	flagYes        bool
	flagPass       string
	flagListStyles bool
	flagListGroups bool
	flagSave       bool
	flagTagline    string
	flagNoTagline  bool
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate & apply a stylised git commit message",
	RunE:  runCommit,
}

func init() {
	commitCmd.Flags().StringVarP(&flagStyle, "style", "s", "", "persona style or 'random'")
	commitCmd.Flags().BoolVarP(&flagRandom, "random", "r", false, "fully random persona")
	commitCmd.Flags().StringVarP(&flagGroup, "group", "g", "", "random persona from this group")
	commitCmd.Flags().StringVarP(&flagMood, "mood", "m", "", "mood or 'random'")
	commitCmd.Flags().StringVarP(&flagLength, "length", "l", "", "short | medium | long")
	commitCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "skip confirmation prompt")
	commitCmd.Flags().StringVarP(&flagPass, "pass-secret", "p", "", "pass secret for GEMINI_API_KEY")
	commitCmd.Flags().BoolVarP(&flagListStyles, "list-styles", "L", false, "list personas & exit")
	commitCmd.Flags().BoolVarP(&flagListGroups, "list-groups", "G", false, "list persona groups & exit")
	commitCmd.Flags().BoolVarP(&flagSave, "save", "S", false, "save current flags as defaults")
	commitCmd.Flags().StringVarP(&flagTagline, "tagline-style", "t", "", "persona for success tagline")
	commitCmd.Flags().BoolVarP(&flagNoTagline, "no-tagline", "T", false, "suppress success tagline")
}

/* ------------------- COMMAND ENTRY ------------------ */

func runCommit(cmd *cobra.Command, _ []string) error {
	if flagListGroups {
		fmt.Println("Available groups:")
		for _, g := range styles.GroupNames() {
			fmt.Printf("  • %s\n", g)
		}
		return nil
	}
	if flagListStyles {
		fmt.Println("Available personas:")
		for _, p := range styles.Personas {
			fmt.Printf("  • %s\n", p)
		}
		return nil
	}

	if _, err := os.Stat(".git"); err != nil {
		return errors.New("❌ not inside a git repository")
	}
	apiKey, err := getAPIKey(flagPass)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	length := pickLength()

	userMsg, err := promptCommitMessage()
	if err != nil {
		if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
			fmt.Println("\n🚫 Aborted.")
			return nil
		}
		return err
	}

	finalMsg, err := confirmFlow(userMsg, length, apiKey)
	if err != nil {
		return err
	}
	if finalMsg == "" {
		fmt.Println("🚫 Aborted.")
		return nil
	}

	if err := gitCommit(finalMsg); err != nil {
		return err
	}
	fmt.Println("🎉 Git commit successful!")

	if !flagNoTagline && viper.GetBool("tagline_enabled") {
		tagPersona := taglinePersona()
		line, _ := gemini.Generate(apiKey, tagPersona, "excited", "short",
			"Celebrate the successful commit with a witty one-liner (≤12 words).")
		fmt.Printf("%s says: %s\n", strings.Title(tagPersona), line)
	}

	if flagSave {
		saveDefaults(cmd, flagTagline, flagGroup, flagMood)
	}
	return nil
}

/* -------------------- HELPERS --------------------- */

func getAPIKey(pass string) (string, error) {
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		return key, nil
	}
	if pass == "" {
		pass = viper.GetString("pass_secret")
	}
	if pass != "" {
		out, err := exec.Command("pass", "show", pass).Output()
		if err == nil {
			if k := strings.TrimSpace(string(out)); k != "" {
				return k, nil
			}
		}
	}
	return "", errors.New("❌ GEMINI_API_KEY not set and no usable pass secret found")
}

func pickStyle() string {
	if flagStyle != "" && strings.ToLower(flagStyle) != "random" {
		return flagStyle
	}
	if flagStyle != "" && strings.ToLower(flagStyle) == "random" {
		return styles.Random()
	}
	if flagGroup != "" {
		if s, ok := styles.RandomFromGroup(flagGroup); ok {
			return s
		}
	}
	if flagRandom {
		return styles.Random()
	}
	if def := viper.GetString("default_character"); def != "" && strings.ToLower(def) != "random" {
		return def
	}
	if grp := viper.GetString("default_group"); grp != "" {
		if s, ok := styles.RandomFromGroup(grp); ok {
			return s
		}
	}
	return styles.Random()
}

func pickMoodOnce() string {
	if flagMood != "" && strings.ToLower(flagMood) != "random" {
		return flagMood
	}
	if def := viper.GetString("default_mood"); def != "" && strings.ToLower(def) != "random" {
		return def
	}
	return styles.RandomMood()
}

func moodIsRandomConfig() bool {
	if flagMood != "" {
		return strings.ToLower(flagMood) == "random"
	}
	return strings.ToLower(viper.GetString("default_mood")) == "random"
}

func pickLength() string {
	if flagLength != "" {
		ll := strings.ToLower(flagLength)
		if ll == "short" || ll == "medium" || ll == "long" {
			return ll
		}
		fmt.Println("⚠️  invalid length; defaulting to medium")
	}
	return viper.GetString("default_length")
}

func promptCommitMessage() (string, error) {
	fmt.Print("💬 Enter your commit message: ")
	r := bufio.NewReader(os.Stdin)
	msg, err := r.ReadString('\n')
	return strings.TrimSpace(msg), err
}

/* ---------------- CONFIRMATION LOOP ---------------- */

func confirmFlow(orig string, length string, apiKey string) (string, error) {
	style := pickStyle()
	mood := pickMoodOnce()
	randomMood := moodIsRandomConfig()
	randomStyle := flagRandom || strings.ToLower(flagStyle) == "random" ||
		(flagGroup != "" && flagStyle == "")

	if flagYes || !viper.GetBool("confirm") {
		return gemini.Generate(apiKey, style, mood, length, orig)
	}

	for {
		if randomStyle {
			style = pickStyle()
		}
		if randomMood {
			mood = styles.RandomMood()
		}

		gen, err := gemini.Generate(apiKey, style, mood, length, orig)
		if err != nil {
			return "", err
		}
		fmt.Printf("\n🧠 Generated commit message (%s, %s, %s):\n\n\"%s\"\n\n",
			style, mood, length, gen)

		// first Y/n prompt
		conf := promptui.Prompt{
			Label:     "✅ Use this message?",
			IsConfirm: true,
			Default:   "Y",
		}
		ans, perr := conf.Run()
		switch perr {
		case promptui.ErrInterrupt, promptui.ErrEOF:
			return "", nil
		case promptui.ErrAbort:
			// typed "n" – fall through to menu
		case nil:
			if ans == "" || strings.ToLower(ans) == "y" {
				return gen, nil
			}
		default:
			return "", perr
		}

		// secondary menu
		menu := promptui.Select{
			Label:        "❓ What next?",
			Items:        []string{"Generate another", "Use my original", "Cancel"},
			HideSelected: true,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "❯ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "  {{ . }}",
			},
		}
		_, act, merr := menu.Run()
		if merr == promptui.ErrInterrupt || merr == promptui.ErrEOF {
			return "", nil
		} else if merr != nil {
			return "", merr
		}

		switch act {
		case "Generate another":
			continue
		case "Use my original":
			return orig, nil
		default: // Cancel
			return "", nil
		}
	}
}

/* ---------------- GIT EXEC & SAVE ---------------- */

func gitCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func taglinePersona() string {
	if flagTagline != "" {
		return flagTagline
	}
	return viper.GetString("tagline_persona")
}

func saveDefaults(cmd *cobra.Command, tagline, group, mood string) {
	changed := false
	if cmd.Flags().Changed("tagline-style") {
		viper.Set("tagline_persona", tagline); changed = true
	}
	if cmd.Flags().Changed("group") {
		viper.Set("default_group", group); changed = true
	}
	if cmd.Flags().Changed("mood") && strings.ToLower(mood) == "random" {
		viper.Set("default_mood", "random"); changed = true
	}

	if changed {
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("⚠️  could not write config: %v\n", err)
		} else {
			fmt.Println("💾 Saved new defaults ✔")
		}
	}
}

