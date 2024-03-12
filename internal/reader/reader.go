package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func Read(path string, ch chan [2]string) error {
	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("opening input file %s: %w", path, err)
		return err
	}
	defer func() {
		_ = f.Close()
		close(ch)
	}()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		terms := strings.Split(scanner.Text(), ";")
		if len(terms) != 2 {
			return errors.New("invalid line")
		}
		ch <- [2]string(terms)
	}

	return nil
}
