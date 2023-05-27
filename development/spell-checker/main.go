package main

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Root directory to check spelling
	rootDir string

	// Personal dictionary
	personalDict string

	// Language
	lang string
)

var rootCmd = &cobra.Command{
	Use:   "spell-checker",
	Short: "Spell checker is a wrapper around aspell",
	Long:  `Spell checker is a wrapper around aspell to check spelling of a files in directory`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := getFilesToCheck(rootDir, `\.(md|tex|txt)$`)
		if err != nil {
			log.Fatalf("Error getting files to check: %s", err)
		}

		for _, file := range files {
			log.Printf("Checking file: %s\n", file)
			words, err := spellCheckFile(file)
			if err != nil {
				log.Fatalf("Error checking file: %s", err)
			}

			if len(words) == 1 {
				log.Printf("No misspelled words found")
				continue
			}

			log.Printf("Found %d misspelled words:\n%s", len(words)-1, strings.Join(words, "\n"))
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", ".", "Root directory to check spelling")
	rootCmd.PersistentFlags().StringVarP(&personalDict, "personal", "p", "", "Personal dictionary")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "en", "Language")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getFilesToCheck(root, ext string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		re := regexp.MustCompile(ext)
		if !re.MatchString(path) {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func spellCheckFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cmd *exec.Cmd
	if personalDict != "" {
		cmd = exec.Command("aspell", "-l", lang, "-p", personalDict, "list")
	} else {
		cmd = exec.Command("aspell", "-l", lang, "list")
	}

	cmd.Stdin = f

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return strings.Split(out.String(), "\n"), nil

}
