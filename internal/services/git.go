package services

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/maycon-jesus/mj-cli/pkg/cmd"
)

var ErrNoMainBranch = fmt.Errorf("no main or master branch found")

type GitService struct {
	logger *slog.Logger
}

func NewGitService() *GitService {
	return &GitService{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func (s *GitService) WithLogger(logger *slog.Logger) *GitService {
	s.logger = logger
	return s
}

func (s *GitService) DetectMainBranch() (string, error) {
	s.logger.Debug("Detecting main branch")

	// Tenta verificar se existe a branch main
	exists, err := s.BranchExists("main")
	if err != nil {
		return "", err
	}

	if exists {
		return "main", nil
	}

	// Tenta verificar se existe a branch master
	exists, err = s.BranchExists("master")
	if err != nil {
		return "", err
	}

	if exists {
		return "master", nil
	}

	// Se nenhuma das duas existir, retorna erro
	s.logger.Debug("No main or master branch found")
	return "", ErrNoMainBranch
}

func (s *GitService) Checkout(branch string) error {
	s.logger.Debug("Checking out branch", "branch", branch)
	err := cmd.RunCommandWithOptions(fmt.Sprintf("git checkout %s", branch), cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to checkout branch", "branch", branch, "error", err)
		return err
	}
	s.logger.Debug("Checked out branch successfully", "branch", branch)
	return nil
}

func (s *GitService) NewBranch(newBranch string) error {
	if exists, err := s.BranchExists(newBranch); err != nil {
		return err
	} else if exists {
		s.logger.Debug("Branch already exists", "branch", newBranch)
		return fmt.Errorf("branch %s already exists", newBranch)
	}

	s.logger.Debug("Creating branch from current HEAD", "branch", newBranch)
	err := cmd.RunCommandWithOptions(fmt.Sprintf("git checkout -b %s", newBranch), cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to create branch", "branch", newBranch, "error", err)
		return err
	}
	s.logger.Debug("Branch created successfully", "branch", newBranch)
	return nil
}

func (s *GitService) Pull() error {
	s.logger.Debug("Pulling latest changes")
	err := cmd.RunCommandWithOptions("git pull", cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to pull", "error", err)
		return err
	}
	s.logger.Debug("Pull completed successfully")
	return nil
}

func (s *GitService) BranchExists(branch string) (bool, error) {
	s.logger.Debug("Checking if branch exists", "branch", branch)
	_, err := cmd.GetCommandOutput(fmt.Sprintf("git rev-parse --verify refs/heads/%s", branch))
	if err != nil {
		s.logger.Debug("Branch does not exist", "branch", branch)
		return false, nil
	}
	s.logger.Debug("Branch exists", "branch", branch)
	return true, nil
}

func (s *GitService) UndoLastCommit() error {
	s.logger.Debug("Undoing last commit")
	err := cmd.RunCommandWithOptions("git reset --soft HEAD~1", cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to undo last commit", "error", err)
		return err
	}
	s.logger.Debug("Last commit undone successfully")
	return nil
}

func (s *GitService) GetCurrentBranch() (string, error) {
	s.logger.Debug("Getting current branch")
	output, err := cmd.GetCommandOutput("git branch --show-current")
	if err != nil {
		s.logger.Debug("Failed to get current branch", "error", err)
		return "", err
	}
	currentBranch := strings.TrimSpace(string(output))
	s.logger.Debug("Current branch retrieved successfully", "branch", currentBranch)
	return currentBranch, nil
}

func (s *GitService) Rebase(baseBranch string) error {
	s.logger.Debug("Rebasing current branch onto base branch", "baseBranch", baseBranch)
	output, err := cmd.GetCommandOutput(fmt.Sprintf("git rebase %s", baseBranch))
	if err != nil {
		s.logger.Debug("Failed to rebase", "baseBranch", baseBranch, "error", err, "output", output)
		return err
	}
	s.logger.Debug("Rebase completed successfully", "baseBranch", baseBranch)
	return nil
}

func (s *GitService) BranchHasNewCommitsFor(baseBranch string) (bool, error) {
	s.logger.Debug("Checking if base branch has new commits not in HEAD", "baseBranch", baseBranch)
	output, err := cmd.GetCommandOutput(fmt.Sprintf("git rev-list --right-only --count HEAD...%s", baseBranch))
	if err != nil {
		s.logger.Debug("Failed to check for new commits", "baseBranch", baseBranch, "error", err, "output", output)
		return false, err
	}
	countStr := strings.TrimSpace(output)
	hasNewCommits := countStr != "0"
	s.logger.Debug("Checked for new commits successfully", "baseBranch", baseBranch, "hasNewCommits", hasNewCommits)
	return hasNewCommits, nil
}

func (s *GitService) HasUncommitedChanges() (bool, error) {
	s.logger.Debug("Checking for uncommited changes")
	output, err := cmd.GetCommandOutput("git status --porcelain")
	if err != nil {
		s.logger.Debug("Failed to check for uncommited changes", "error", err)
		return false, err
	}
	hasChanges := strings.TrimSpace(output) != ""
	s.logger.Debug("Checked for uncommited changes successfully", "hasUncommitedChanges", hasChanges)
	return hasChanges, nil
}

func (s *GitService) StashPush(stashName string) error {
	s.logger.Debug("Stashing changes", "stashName", stashName)
	err := cmd.RunCommandWithOptions(fmt.Sprintf("git stash push -m \"%s\"", stashName), cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to stash changes", "stashName", stashName, "error", err)
		return err
	}
	s.logger.Debug("Changes stashed successfully", "stashName", stashName)
	return nil
}

func (s *GitService) StashPop() error {
	s.logger.Debug("Popping latest stash")
	err := cmd.RunCommandWithOptions("git stash pop", cmd.CommandOptions{})
	if err != nil {
		s.logger.Debug("Failed to pop latest stash", "error", err)
		return err
	}
	s.logger.Debug("Latest stash popped successfully")
	return nil
}
