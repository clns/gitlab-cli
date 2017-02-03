package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func SaveViperConfig() error {
	filename := viper.ConfigFileUsed()
	if filename == "" {
		hdir, err := homedir.Dir()
		if err != nil {
			return err
		}
		filename = filepath.Join(hdir, configName+".yml")
	}
	f, err := os.Create(filename)
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
