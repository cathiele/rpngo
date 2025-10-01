package startup

import (
	"errors"
	"fmt"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os"
	"path/filepath"
)

const defaultConfigFile = `
# Create and layout windows
'g1' w.new.group
'g1' w.columns
'i' 'g1' w.move.end
'g2' w.new.group
'g2' 25 w.weight
'g2' w.columns
's' w.new.stack
's' 'g2' w.move.end
'v' w.new.var
'v' 'g2' w.move.end

# Plot defaults
'p' plot.win=
'$plot.win w.new.plot $plot.win "g1" w.move.end' plot.init=
` + commonStartup

const configName = ".rpngo"

// OSStartup is startup logic when running in an os context (e.g. with a
// filesystem available)
func OSStartup(r *rpn.RPN) error {
	configPath, err := genConfigPath()
	if err != nil {
		return err
	}
	s, err := loadOrCreateConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("while loading %s: %w", configPath, err)
	}
	fields, err := parse.Fields(s)
	if err != nil {
		return fmt.Errorf("while parsing %s: %w", configPath, err)
	}
	if err := r.Exec(fields); err != nil {
		return fmt.Errorf("while executing commands in %s: %w", configPath, err)
	}
	return nil
}

func genConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configName), nil
}

func loadOrCreateConfigFile(configPath string) (string, error) {
	s, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		s, err = createConfigFile(configPath)
	}
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func createConfigFile(configPath string) ([]byte, error) {
	s := []byte(defaultConfigFile)
	if err := os.WriteFile(configPath, s, 0644); err != nil {
		return nil, err
	}
	return s, nil
}
