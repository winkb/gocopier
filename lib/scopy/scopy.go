package scopy

import "github.com/atotto/clipboard"

func CopyText(text string) error {
	var err = clipboard.WriteAll(text)
	return err
}
