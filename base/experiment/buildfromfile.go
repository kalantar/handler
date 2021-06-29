package experiment

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// FromFile builds an experiment from a yaml file.
func (b *Builder) FromFile(filePath string) *Builder {
	var err error
	var data []byte
	if data, err = ioutil.ReadFile(filePath); err == nil {
		exp := &Experiment{}
		if err = yaml.Unmarshal(data, exp); err == nil {
			actions, _ := json.MarshalIndent(exp.Spec.Strategy.Actions, "", "  ")
			log.Trace(string(actions))
			b.exp = exp
			return b
		}
		log.Error(err)
	}
	log.Error(err)
	b.err = errors.New("cannot build experiment from file")
	return b
}
