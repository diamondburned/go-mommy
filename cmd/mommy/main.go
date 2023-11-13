package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"libdb.so/go-mommy"
)

// mommy variables
var (
	moods             []string
	emotes            []string
	pronouns          []string
	roles             []string
	affectionateTerms []string
	denigratingTerms  []string
	parts             []string
)

// cmd variables
var (
	responsesFile = ""
	stylize       = false
	nsfw          = false
	seed          = int64(time.Now().UnixNano())
)

func defaultValues(key mommy.VariableKey) []string {
	return mommy.DefaultResponses.Vars[key].Defaults
}

func init() {
	log.SetFlags(0)

	pflag.BoolVarP(&stylize, "stylize", "s", stylize,
		"stylize the output using ANSI escape codes")
	pflag.BoolVarP(&nsfw, "nsfw", "n", nsfw,
		"show NSFW flags and include NSFW responses")
	pflag.Int64VarP(&seed, "seed", "S", seed,
		"seed for the random number generator")
	pflag.StringVarP(&responsesFile, "responses-file", "f", responsesFile,
		"responses.json file (default: bundled responses.json)")

	// Determine if our flags should be --mommy-* or --daddy-*.
	role := "mommy"
	if len(roles) == 1 {
		// Use the first role.
		role = roles[0]
	} else {
		// Guess from os.Args[0], falling back to "mommy".
		if filepath.Base(os.Args[0]) == "daddy" {
			role = "daddy"
		}
	}

	replaceRole := func(s string) string { return strings.ReplaceAll(s, "<>", role) }
	prefix := role + "s" // possessive form

	pflag.StringSliceVar(&affectionateTerms, prefix+"-little", affectionateTerms,
		replaceRole("what to call you~"))
	pflag.StringSliceVar(&pronouns, prefix+"-pronouns", pronouns,
		replaceRole("what pronouns <> will use for themself~"))
	pflag.StringSliceVar(&roles, prefix+"-roles", roles,
		replaceRole("what role <> will have~"))
	pflag.StringSliceVar(&emotes, prefix+"-emotes", emotes,
		replaceRole("what emotes <> will have~"))

	pflag.StringSliceVar(&moods, prefix+"-moods", moods,
		replaceRole("how kinky <> will be~ (nsfw)"))
	pflag.StringSliceVar(&parts, prefix+"-parts", parts,
		replaceRole("what part of <> you should crave~ (nsfw)"))
	pflag.StringSliceVar(&denigratingTerms, prefix+"-fucking", denigratingTerms,
		replaceRole("what to call <>'s pet~ (nsfw)"))

	pflag.Usage = func() {
		if !nsfw {
			pflag.CommandLine.MarkHidden(prefix + "-moods")
			pflag.CommandLine.MarkHidden(prefix + "-parts")
			pflag.CommandLine.MarkHidden(prefix + "-fucking")
		}

		log.Println("Usage:")
		log.Println("  mommy [options] <response-type>")
		log.Println("  daddy [options] <response-type>")
		log.Println("")
		log.Println("  <response-type> := <positive> | <negative>")
		log.Println("  <positive>      := positive | + | 0")
		log.Println("  <negative>      := negative | - | 1")
		log.Println("")

		log.Println("Examples:")
		log.Println("  mommy --mommys-little girl positive")
		log.Println("  daddy --daddys-little boy --daddys-pronouns his -")
		log.Println("")

		log.Println("Flags:")
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()

	var responseType mommy.ResponseType
	switch arg := pflag.Arg(0); arg {
	case "positive", "+", "0":
		responseType = mommy.PositiveResponse
	case "negative", "-", "1":
		responseType = mommy.NegativeResponse
	default:
		fatal(fmt.Sprintf("unknown response type %q, want positive|negative", arg))
	}

	config := mommy.DefaultResponses
	if responsesFile != "" {
		b, err := os.ReadFile(responsesFile)
		if err != nil {
			fatal("cannot read responses file:", err)
		}
		config, err = mommy.UnmarshalResponses(b)
		if err != nil {
			fatal("cannot parse responses file:", err)
		}
	}

	overrideDefaults(&config, mommy.VariableMood, moods)
	overrideDefaults(&config, mommy.VariableEmote, emotes)
	overrideDefaults(&config, mommy.VariablePronoun, pronouns)
	overrideDefaults(&config, mommy.VariableRole, roles)
	overrideDefaults(&config, mommy.VariableAffectionateTerm, affectionateTerms)
	overrideDefaults(&config, mommy.VariableDenigratingTerm, denigratingTerms)
	overrideDefaults(&config, mommy.VariablePart, parts)

	if !nsfw && isNSFW(config) {
		fatal("cannot generate NSFW content without --nsfw")
	}

	rng := rand.New(rand.NewSource(seed))

	gen, err := mommy.NewGeneratorWithRandom(config, rng)
	if err != nil {
		fatal(err)
	}

	res, err := gen.Generate(responseType, nil)
	if err != nil {
		fatal(err)
	}

	res = stylizeResponse(res)
	fmt.Println(res)
}

func overrideDefaults(cfg *mommy.Responses, key mommy.VariableKey, defaults []string) {
	if len(defaults) > 0 {
		v := cfg.Vars[key]
		v.Defaults = defaults
		cfg.Vars[key] = v
	}
}

func fatal(v ...any) {
	gen, _ := mommy.NewGeneratorWithRandom(mommy.DefaultResponses, rand.New(rand.NewSource(seed)))
	res, _ := gen.Generate(mommy.NegativeResponse, mommy.Overrides{
		mommy.VariableMood: mommy.Chill,
	})
	log.Println(res)
	log.Fatalln(v...)
}

func stylizeResponse(res string) string {
	// The responses only ever have *italic* on its own line.
	// This simplifies the logic a lot.
	lines := strings.Split(res, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "*") && strings.HasSuffix(line, "*") {
			lines[i] = "\033[3m" + line + "\033[0m"
		}
	}
	return strings.Join(lines, "\n")
}

func isNSFW(cfg mommy.Responses) bool {
	moods := cfg.Vars[mommy.VariableMood].Defaults
	return slices.ContainsFunc(moods, func(s string) bool { return s != mommy.Chill })
}
