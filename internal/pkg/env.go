package pkg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type EnvValue struct {
	pairs  []string
	latest map[string]int
}

func (ev *EnvValue) String() string {
	return fmt.Sprintf("%q", ev.pairs)
}

var (
	ErrEmptyKey = errors.New("empty key")
)

func (ev *EnvValue) Set(s string) error {
	k, v, found := strings.Cut(s, "=")
	if k == "" {
		return ErrEmptyKey
	}
	// -e HOSTNAME
	if !found {
		var ok bool

		v, ok = os.LookupEnv(k)
		if !ok {
			return nil
		}
	}

	if ev.latest == nil {
		ev.latest = make(map[string]int)
	}

	if i, ok := ev.latest[k]; ok {
		ev.pairs[i] = k + "=" + v
	} else {
		ev.pairs = append(ev.pairs, k+"="+v)
		ev.latest[k] = len(ev.pairs) - 1
	}

	return nil
}

func (ev *EnvValue) Env() []string {
	c := make([]string, len(ev.pairs))
	copy(c, ev.pairs)

	return c
}
