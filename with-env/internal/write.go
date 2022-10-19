package internal

import (
	"bufio"
	"io"
	"strings"
)

const (
	redactComment = "# with-env redact"
	redacted      = "<redacted>\n"
)

func RedactWrite(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		// Redact lines following a redact comment.
		if strings.TrimSpace(text) == redactComment {
			ok, err := redactAtComment(w, text, scanner)
			if err != nil {
				return err
			}
			if !ok {
				// Skip to checking scanner.Err
				break
			}
			continue
		}

		// Auto-redact anything where the key includes "password".
		if strings.Contains(strings.ToLower(strings.TrimSpace(text)), "password") {
			wrote, err := autoRedactPasswords(w, text, scanner)
			if err != nil {
				return err
			}
			if wrote {
				continue
			}
		}

		// No applied redaction.
		if err := writeString(w, text+"\n"); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func writeString(w io.Writer, v string) error {
	// io.Write{}.Write _must_ return an error is n is less than len(v).
	// To that end, I don't care about n.
	_, err := io.WriteString(w, v)
	return err
}

func redactAtComment(w io.Writer, text string, scanner *bufio.Scanner) (bool, error) {
	if err := writeString(w, text+"\n"); err != nil {
		return false, err
	}

	if !scanner.Scan() {
		return false, nil
	}

	var err error
	fields := strings.Split(strings.TrimSpace(scanner.Text()), "=")
	if len(fields) != 2 {
		err = writeString(w, redacted)
	} else {
		fields[1] = redacted
		err = writeString(w, strings.Join(fields, "="))
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func autoRedactPasswords(w io.Writer, text string, scanner *bufio.Scanner) (bool, error) {
	fields := strings.Split(strings.TrimSpace(text), "=")
	if len(fields) == 2 {
		fields[1] = redacted
		if err := writeString(w, strings.Join(fields, "=")); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
