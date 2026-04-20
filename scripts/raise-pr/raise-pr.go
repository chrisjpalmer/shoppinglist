package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/*
commit 46f658d394fd02d7e2b0dd175f8f316ad9b4e645
Merge: e298ed2 c7958a7
Author: Chris Palmer <chrisjpalmer6@gmail.com>
Commit: GitHub <noreply@github.com>

    Merge pull request #65 from chrisjpalmer/introduce-claude-settings

    🔒 Introduce claude permission settings

commit 304fa7151caabe4fc60ffd5d1c7940b5f16c4f3f (origin/reorganise-docs, reorganise-docs)
Author: Chris Palmer <chrisjpalmer6@gmail.com>
Commit: Chris Palmer <chrisjpalmer6@gmail.com>

    🚚 Reorganise docs directory

    Reorganise the docs directory by introducing a dedicated images
    subdirectory for images used by markdown files.
*/

func main() {
	if err := raisePr(); err != nil {
		log.Fatal(err)
	}
}

func raisePr() error {
	subject, body, err := lastCommitMessage()
	if err != nil {
		return fmt.Errorf("error getting last commit message: %w", err)
	}

	if err := raiseGhPr(subject, body); err != nil {
		return fmt.Errorf("failed to raise the github pr: %w", err)
	}

	return nil
}

func lastCommitMessage() (string, string, error) {
	out, err := shellExec("git", "log", "-n", "1", "--format=full")
	if err != nil {
		return "", "", fmt.Errorf("error executing git log: %w", err)
	}

	lines := strings.Split(out, "\n")

	if len(lines) < 5 {
		return "", "", fmt.Errorf("expected at least 5 lines for the commit info but only have %d", len(lines))
	}

	if strings.HasPrefix(lines[1], "Merge:") {
		return "", "", fmt.Errorf("current commit is a merge commit, can't use for raising a PR")
	}

	subject := strings.TrimSpace(lines[4])

	if len(lines) < 6 {
		return subject, "", nil
	}

	bodyLines := lines[5:]

	trimLines(bodyLines)

	body := strings.Join(bodyLines, " ")

	return subject, body, nil
}

func raiseGhPr(subject, body string) error {
	out, err := shellExec("gh", "pr", "create", "--title", subject, "--body", body)

	if err != nil {
		return fmt.Errorf("error while creating github pr: %w", err)
	}

	fmt.Println(out)

	return nil
}

func trimLines(lines []string) {
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
}

func shellExec(cmdStr string, args ...string) (string, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command(cmdStr, args...)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
