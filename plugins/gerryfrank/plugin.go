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

	numWords := rand.Intn(15) + 20
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

	generatedLine := c.Generate(numWords)

	bot.Send(person, generatedLine)
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
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)

	// start with a random word from the chain
	candidates := make([]string, 0, len(c.chain))
	for k := range c.chain {
		candidates = append(candidates, k)
	}

	p.Shift(candidates[rand.Intn(len(candidates))])

	var words []string
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
	return strings.Join(words, " ")
}

func getRandomPerson() string {
	persons := []string{"pat", "gerry", "frank", "audience"}
	return persons[rand.Intn(len(persons))]
}
