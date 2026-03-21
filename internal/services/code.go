package services

import (
	"log/slog"
	"path"

	"github.com/maycon-jesus/mj-cli/pkg/filesystem"
)

type CodeService struct {
	logger *slog.Logger
}

func NewCodeService() *CodeService {
	return &CodeService{
		logger: slog.New(slog.NewTextHandler(nil, nil)),
	}
}

func (s *CodeService) WithLogger(logger *slog.Logger) *CodeService {
	s.logger = logger
	return s
}

func (s *CodeService) GetVersion(projectPath string) (string, error) {
	s.logger.Debug("Getting current version from code")
	versionPath := path.Join(projectPath, "VERSION")
	return filesystem.GetFileContent(versionPath)
}

func (s *CodeService) SetVersion(projectPath string, version string) error {
	s.logger.Debug("Setting new version in code", "version", version)
	versionPath := path.Join(projectPath, "VERSION")
	return filesystem.WriteFile(versionPath, version)
}
