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
#"
#'g' w.new.group
#'g' w.columns
#'i' 'g' w.move.end
#'i' 25 w.weight
#'g2' w.new.group
#'g2' 25 w.weight
#'g2' w.columns
#'g2' .wtarget=
#'s' w.new.stack
#'v' w.new.var
#'g' .wtarget=
#" .init=
#
#@.init

'w.reset @.init' .f1=
'w.reset "root" w.columns "i" 30 w.weight' .f2=
'w.reset "root" w.columns "v" w.new.var "v" "showdot" true w.setp' .f3=

'time t1= 0 x= "$x 1 + x= $x 3000000 <" for time $t1 - 3000000 1> /' benchmark=

# Plot defaults
'p' .plotwin=
'$.plotwin w.new.plot' .plotinit=
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
	fields := make([]string, 256)
	fields, err = parse.Fields(s, fields)
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
