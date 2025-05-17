package styles

import (
	"math/rand"
	"sort"
	"strings"
)

/* ----------------------------------- */
/*    GROUP DEFINITIONS & PERSONAS     */
/* ----------------------------------- */

// Groups maps lowercase group names to their personas.
var Groups = map[string][]string{
	"trailer_park_boys": {
		"ricky", "julian", "bubbles", "jim lahey", "randy (tpb)", "j-roc",
		"cory", "trevor", "lucy (tpb)", "sarah (tpb)", "ray (tpb)", "trinity",
		"sam losco", "conky", "barb lahey", "jacob collins", "phil collins (tpb)",
		"george green", "cyrus", "t-bone", "don (liquor store)", "terry (samsquamptch hunter)",
	},
	"cartoons": {
		"yoda", "homer simpson", "rick sanchez", "bender", "spongebob",
		"stewie griffin", "kermit the frog", "gollum", "glados",
	},
	"politicians": {
		"donald trump", "barack obama", "bernie sanders",
	},
	"celebrities": {
		"elon musk", "morgan freeman", "keanu reeves", "ron burgundy",
		"deadpool", "tony stark", "jack sparrow", "captain kirk",
	},
	"literary": {
		"shakespeare", "gandalf", "mark twain",
	},
	"misc": {
		"doge", "strong bad", "ace ventura", "monty python narrator",
		"ivar aasen", // special case stays
	},
}

/* Flatten everything into Personas for backward-compat */
var Personas []string

func init() {
	seen := make(map[string]bool)
	for _, list := range Groups {
		for _, p := range list {
			l := strings.ToLower(p)
			if !seen[l] {
				Personas = append(Personas, p)
				seen[l] = true
			}
		}
	}
	sort.Strings(Personas)
}

/* ----------------------------------- */
/*            HELPER FUNCS             */
/* ----------------------------------- */

func Random() string {
	return Personas[rand.Intn(len(Personas))]
}

func RandomFromGroup(group string) (string, bool) {
	if list, ok := Groups[strings.ToLower(group)]; ok && len(list) > 0 {
		return list[rand.Intn(len(list))], true
	}
	return "", false
}

func GroupNames() []string {
	var names []string
	for k := range Groups {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}


/* ----------------------- moods ------------------------ */

var Moods = []string{
	"playful", "sarcastic", "enthusiastic", "melancholic", "dramatic",
	"epic", "witty", "mysterious", "angry", "poetic",
}

func RandomMood() string {
	return Moods[rand.Intn(len(Moods))]
}

