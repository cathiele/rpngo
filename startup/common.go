package startup

import (
	"fmt"
	"mattwach/rpngo/elog"
	"mattwach/rpngo/fileops"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os"
	"path/filepath"
)

const commonStartup = `
# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
# (-b +/- sqrt(b*b - 4*a*c)) / (2 * a)
{$2 * 4 * $1 sq - neg sqrt 1> neg $0 $2 - $3 2 * / 3< + 1> 2 * /} quad=

{0 {+ ssize 1 >} for} sum=
{ssize n< 0 {+ ssize 1 >} for n> /} mean=
{$0 {min ssize 1 >} for} min=
{$0 {max ssize 1 >} for} max=

histl
hists
`

const configName = ".rpngo"

// Startup tries to load .rpngo and tries to create a default
// file if one can not be loaded.
func Startup(r *rpn.RPN, fs fileops.FileOpsDriver) error {
	configPath, err := genConfigPath()
	if err != nil {
		return err
	}
	s := loadOrCreateConfigFile(fs, configPath)
	err = parse.Fields(s, r.Exec)
	if err != nil {
		return fmt.Errorf("while parsing %s: %w", configPath, err)
	}
	return nil
}

func genConfigPath() (string, error) {
	home, err := fileops.HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configName), nil
}

func loadOrCreateConfigFile(fs fileops.FileOpsDriver, configPath string) string {
	s, err := fs.ReadFile(configPath)
	if err != nil {
		elog.Print("while loading config ", configPath, ": ", err.Error())
		s = createConfigFile(fs, configPath)
	}
	return string(s)
}

func createConfigFile(fs fileops.FileOpsDriver, configPath string) []byte {
	s := []byte(defaultConfig)
	if err := os.WriteFile(configPath, s, 0644); err != nil {
		elog.Print("while saving config ", configPath, ": ", err.Error())
		s = createConfigFile(fs, configPath)
	}
	return s
}
