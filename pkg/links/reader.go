package links

import (
	"bufio"
	"os"
	"strings"
)

// ReadLinksList provides buffered reading file line by line
// It is suitable for big files
func ReadLinksList(filename string, bufSize uint) (links chan string, err error) {
	links = make(chan string, bufSize)

	file, err := os.Open(filename)
	if err != nil {
		return links, err
	}

	scanner := bufio.NewReader(file)
	go func() {
		for {
			line, _, err := scanner.ReadLine()
			if err != nil {
				break
			}
			str := strings.TrimSpace(string(line))
			if len(str) == 0 {
				continue
			}
			links <- str
		}
		file.Close()
		close(links)
	}()

	return links, nil
}
