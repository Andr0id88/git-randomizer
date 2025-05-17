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
	"cartoons": {
		"yoda", "homer simpson", "rick sanchez", "bender", "spongebob",
		"stewie griffin", "kermit the frog", "gollum", "glados", "eric cartman",
		"beavis", "butt-head", "patrick star", "velma", "meatwad",
	},
	"politicians": {
		"donald trump", "barack obama", "bernie sanders", "kim jong-un",
		"george w. bush", "hillary clinton", "joe biden", "justin trudeau",
		"boris johnson", "arnold schwarzenegger (gov)", "volodymyr zelenskyy",
	},
	"celebrities": {
		"elon musk", "morgan freeman", "keanu reeves", "ron burgundy",
		"deadpool", "tony stark", "jack sparrow", "captain kirk",
		"gordon ramsay", "michael jackson", "nicolas cage", "kanye west",
	},
	"literary": {
		"shakespeare", "gandalf", "mark twain", "edgar allan poe",
		"oscar wilde", "dr. seuss", "j.r.r. tolkien",
	},
	"misc": {
		"doge", "strong bad", "ace ventura", "monty python narrator", "pickle rick",
		"ivar aasen", "bob ross", "mr. bean", "clippy", "glitch gremlin",
		"the intern", "404 bot", "your manager", "chatgpt hallucination",
	},
	"action_heroes": {
		"john wick", "arnold schwarzenegger", "sylvester stallone", "chuck norris",
		"macho man randy savage", "the rock",
	},
	"tech_legends": {
		"steve jobs", "bill gates", "linus torvalds", "mark zuckerberg",
		"richard stallman", "grug brained dev",
	},
	"musicians": {
		"freddie mercury", "bob dylan", "elvis presley", "prince",
		"taylor swift", "david bowie",
	},
	"sci_fi": {
		"neo", "morpheus", "spock", "the borg",
		"darth vader", "emperor palpatine", "yoda (sith edition)",
	},
	"actors": {
		"jack nicholson", "samuel l. jackson", "christopher walken", "jeff goldblum",
		"danny devito", "bill murray", "will ferrell", "owen wilson",
	},
	"internet_legends": {
		"shrek", "keyboard cat", "nyan cat", "pepe the frog",
		"bad luck brian", "unhelpful high school teacher", "overly attached girlfriend",
	},
	"supervillains": {
		"joker", "lex luthor", "thanos", "dr. evil",
		"megamind", "gru", "lord farquaad", "sylar",
	},
	"philosophers": {
		"socrates", "plato", "aristotle", "friedrich nietzsche",
		"karl marx", "rene descartes", "simone de beauvoir", "confucius",
	},
	"conspiracy_theorists": {
		"alex jones", "fox mulder", "david icke", "qanon shaman",
		"flat earth guy", "chemtrail lady", "ancient aliens guy",
	},
	"game_characters": {
		"mario", "luigi", "kratos", "solid snake",
		"gordon freeman", "lara croft", "master chief", "pac-man",
	},
	"robots/ai": {
		"hal 9000", "skynet", "data", "optimus prime",
		"wall-e", "marvin the paranoid android", "roboCop",
	},
	"comedians": {
		"george carlin", "mitch hedberg", "bo burnham", "robin williams",
		"john mulaney", "ricky gervais", "dave chappelle",
	},
	"rappers": {
		"snoop dogg", "eminem", "dr. dre", "tupac",
		"notorious b.i.g.", "kendrick lamar", "ice cube", "missy elliott",
	},
	"rockstars": {
		"ozzy osbourne", "kurt cobain", "axl rose", "jim morrison",
		"jimi hendrix", "mick jagger", "bono", "angus young",
	},
	"trailer_park_boys": {
		"ricky", "julian", "bubbles", "jim lahey", "randy (tpb)", "conky",
		"j-rock", "cyrus", "sam losco", "ray (tpb)", "phil collins (tpb)",
		"trevor", "cory",
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
	"chaotic", "apathetic", "delusional", "bitter", "eccentric",
	"confused", "heroic", "unhinged", "gremlin", "sassy",
	"doomcore", "overconfident", "tragic", "existential",

}

func RandomMood() string {
	return Moods[rand.Intn(len(Moods))]
}

