package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// config log
	log.SetFlags(log.Lshortfile)

	// read config
	viper.SetConfigType("yaml")
	viper.SetConfigName(".todosconfig")
	viper.AddConfigPath("$HOME")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("fatal:", err)
	}

	// set flags
	rootCmd.Flags().StringVarP(&ignore, "ignore", "i", "", "ignore pattern")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
}

var (
	_            log.Logger
	wd           string
	verbose      bool
	ignore       string
	ignoreRe     *regexp.Regexp
	extensions   string
	extensionsRe *regexp.Regexp
	tags         []tag
	ignored      []string
	matched      []string
)

type tag struct {
	file    string
	tag     string
	line    int
	message string
}

func (t tag) Nil() bool {
	return len(t.tag) > 0
}

func Execute() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:  "todos",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		readExtensionsConfig()
		readIgnoreConfig()

		// iterate args paths
		for _, arg := range args {
			fp := path.Join(wd, arg)
			fi, err := os.Lstat(fp)
			if err != nil {
				continue
			}

			if fi.IsDir() {
				if err := iterateDir(fp); err != nil {
					log.Fatal(err)
				}
			} else {
				if !extensionsRe.MatchString(fi.Name()) {
					continue
				}

				if err := parseFile(fp); err != nil {
					log.Fatal(err)
				}
			}
		}

		printTags(tags)
	},
}

func readExtensionsConfig() {
	// get extensions from config
	extensions = strings.Join(viper.GetStringSlice("extensions"), "|")
	extensionsRe = regexp.MustCompile(fmt.Sprintf(`.*\.(%v)$`, extensions))

	if verbose {
		extStr := strings.Join(viper.GetStringSlice("extensions"), ", ")
		fmt.Println("Looking for these extensions:", extStr)
	}
}

func readIgnoreConfig() {
	// get ignore from flag
	ignore = strings.Replace(ignore, "*", ".*", -1)
	ignoreRe = regexp.MustCompile(ignore)

}

func iterateDir(d string) error {
	fis, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {
		name := fi.Name()
		fp := path.Join(d, name)

		if len(ignore) > 0 && ignoreRe.MatchString(fp) {
			ignored = append(ignored, fp)
			continue
		}

		if fi.IsDir() {
			iterateDir(fp)
		} else {
			if extensionsRe.MatchString(name) {
				tags = append(tags, parseFile(fp)...)

				matched = append(matched, fp)
			}
		}
	}

	return nil
}

func parseFile(fp string) (_tags []tag) {
	switch path.Ext(fp) {
	case ".go":
		_tags = ParseGoFile(fp)
	case ".js":
		_tags = ParseGoFile(fp)
	}

	return _tags
}

func printTags(_tags []tag) {
	for _, t := range _tags {
		fmt.Printf(".%v:%v: %v: %v\n", t.file, t.line, t.tag, t.message)
	}
}
