package services

import (
	"fmt"
	"io"
	"log/slog"

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
