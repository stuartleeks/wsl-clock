package wsl

import (
	"fmt"
	"regexp"
	"strings"
)

/*
  NAME                   STATE           VERSION
* Ubuntusl               Running         2
  Ubuntu-18.04           Running         2
  Ubuntu-20.04           Running         2
  golang                 Stopped         2
  dotnet-test            Stopped         2
  docker-desktop-data    Running         2
  docker-desktop         Running         2
*/

type WslDistro struct {
	IsDefault bool
	Name      string
	State     string
	Version   string
}

var listHeaderRegex *regexp.Regexp = regexp.MustCompile("[A-Z]+")

func parseDistroOutput(listOutput string) ([]WslDistro, error) {

	lines := strings.Split(listOutput, "\n")

	header := lines[0]
	lines = lines[1:]

	headerMatches := listHeaderRegex.FindAllIndex([]byte(header), -1)

	checkHeader := func(index int, expectedValue string) error {
		value := strings.TrimSpace(header[headerMatches[index][0]:headerMatches[index][1]])
		if value != expectedValue {
			return fmt.Errorf("Expected header %d to be %q, got %q", index, expectedValue, value)
		}
		return nil
	}
	if len(headerMatches) != 3 {
		return []WslDistro{}, fmt.Errorf("Unexpected headers (expected 3, got %d): %q", len(headerMatches), header)
	}
	if err := checkHeader(0, "NAME"); err != nil {
		return []WslDistro{}, err
	}
	if err := checkHeader(1, "STATE"); err != nil {
		return []WslDistro{}, err
	}
	if err := checkHeader(2, "VERSION"); err != nil {
		return []WslDistro{}, err
	}

	getValue := func(line string, start int, end int) string {
		return strings.TrimSpace(line[start:end])
	}
	distros := []WslDistro{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		distro := WslDistro{}
		if getValue(line, 0, headerMatches[0][0]) == "*" {
			distro.IsDefault = true
		}
		distro.Name = getValue(line, headerMatches[0][0], headerMatches[1][0])
		distro.State = getValue(line, headerMatches[1][0], headerMatches[2][0])
		distro.Version = getValue(line, headerMatches[2][0], len(line))
		distros = append(distros, distro)
	}

	return distros, nil
}
