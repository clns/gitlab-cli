package cmd

import (
	"os"

	"fmt"

	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func SaveViperConfig() error {
	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		return err
	}
	defer f.Close()

	all := viper.AllSettings()
	for k, _ := range all {
		if strings.HasPrefix(k, "_") {
			delete(all, k)
		}
	}
	b, err := yaml.Marshal(all)
	if err != nil {
		return fmt.Errorf("Panic while encoding into YAML format.")
	}
	if _, err := f.WriteString(string(b)); err != nil {
		return err
	}
	return nil
}
