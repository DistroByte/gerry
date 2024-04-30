package gerryfrank

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"git.dbyte.xyz/distro/gerry/shared"
)

func Setup(bot shared.Bot) {
	bot.Register("gerryfrank", Call)
}

// Call generates a random sentence in the style of Gerry, Frank, or the audience.
func Call(bot shared.Bot, context shared.MessageContext, arguments []string, config shared.Config) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	person := getRandomPerson()

	numWords := rand.Intn(15) + 10
	c := NewChain(1)
	// build the chain
	// c.Build(strings.NewReader(TextSources[person]))

	// or load it from a file
	file, err := os.Open(fmt.Sprintf("%s/gerryfrank/%s.json", config.PluginPath, person))
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&c.chain)
	if err != nil {
		return err
	}

	generatedLine := c.Generate(numWords, arguments)

	if generatedLine == "" {
		return nil
	}

	finalString := fmt.Sprintf("%s: %s", person, generatedLine)

	bot.Send(context.Source, context.Target, finalString)
	return nil
}

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
}

func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int, inputWords []string) string {
	p := make(Prefix, c.prefixLen)
	var words []string

	if len(inputWords) > 0 {
		p.Shift(inputWords[rand.Intn(len(inputWords))])

		if _, ok := c.chain[p.String()]; !ok {
			for k := range c.chain {
				p.Shift(k)
				break
			}
		}

	} else {
		for k := range c.chain {
			p.Shift(k)
			break
		}
	}

	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}

		// if the last word of the prefix is a sentence terminator, start a new sentence
		if strings.HasSuffix(p[len(p)-1], ".") || strings.HasSuffix(p[len(p)-1], "!") || strings.HasSuffix(p[len(p)-1], "?") && i > n/2 {
			break
		}

		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}

	if len(words) < 2 {
		return ""
	}

	return strings.Join(words, " ")
}

func getRandomPerson() string {
	persons := []string{"pat", "gerry", "frank", "audience"}
	return persons[rand.Intn(len(persons))]
}
