package services

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type SemverVersion struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

func (v *SemverVersion) String() string {
	version := strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
	if v.PreRelease != "" {
		version += "-" + v.PreRelease
	}
	if v.Build != "" {
		version += "+" + v.Build
	}
	return version
}

func (v *SemverVersion) IncrementMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
	v.ResetBuild()
	v.ResetPreRelease()
}

func (v *SemverVersion) IncrementMinor() {
	v.Minor++
	v.Patch = 0
	v.ResetBuild()
	v.ResetPreRelease()
}

func (v *SemverVersion) IncrementPatch() {
	v.Patch++
	v.ResetBuild()
	v.ResetPreRelease()
}

func (v *SemverVersion) ResetPreRelease() {
	v.PreRelease = ""
}

func (v *SemverVersion) ResetBuild() {
	v.Build = ""
}

var semverRegex = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

var ErrInvalidSemver = errors.New("invalid semantic version")

type SemverService struct {
}

func NewSemverService() *SemverService {
	return &SemverService{}
}

func (s *SemverService) IsValidVersion(version string) bool {
	return semverRegex.MatchString(version)
}

func (s *SemverService) ParseVersion(version string) (*SemverVersion, error) {
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, ErrInvalidSemver
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, ErrInvalidSemver
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, ErrInvalidSemver
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, ErrInvalidSemver
	}

	return &SemverVersion{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: matches[4],
		Build:      matches[5],
	}, nil
}

func (s *SemverService) CompareVersions(v1, v2 string) (int, error) {
	ver1, err := s.ParseVersion(v1)
	if err != nil {
		return 0, err
	}

	ver2, err := s.ParseVersion(v2)
	if err != nil {
		return 0, err
	}

	if ver1.Major != ver2.Major {
		if ver1.Major > ver2.Major {
			return 1, nil
		}
		return -1, nil
	}

	if ver1.Minor != ver2.Minor {
		if ver1.Minor > ver2.Minor {
			return 1, nil
		}
		return -1, nil
	}

	if ver1.Patch != ver2.Patch {
		if ver1.Patch > ver2.Patch {
			return 1, nil
		}
		return -1, nil
	}

	if ver1.PreRelease == "" && ver2.PreRelease != "" {
		return 1, nil
	} else if ver1.PreRelease != "" && ver2.PreRelease == "" {
		return -1, nil
	}

	return s.comparePreRelease(ver1.PreRelease, ver2.PreRelease), nil
}

func (s *SemverService) getPreReleaseIdentifiers(pr string) []string {
	if pr == "" {
		return []string{}
	}
	return strings.Split(pr, ".")
}

func (s *SemverService) comparePreRelease(pr1, pr2 string) int {
	ids1 := s.getPreReleaseIdentifiers(pr1)
	ids2 := s.getPreReleaseIdentifiers(pr2)

	for i := 0; i < len(ids1) && i < len(ids2); i++ {
		if ids1[i] == ids2[i] {
			continue
		}

		num1, err1 := strconv.Atoi(ids1[i])
		num2, err2 := strconv.Atoi(ids2[i])

		switch {
		case err1 == nil && err2 == nil:
			if num1 > num2 {
				return 1
			} else if num1 < num2 {
				return -1
			}
		case err1 == nil && err2 != nil: // Numeric identifiers have lower precedence than non-numeric
			return -1
		case err1 != nil && err2 == nil: // Non-numeric identifiers have higher precedence than numeric
			return 1
		default: // Both are non-numeric, compare lexically
			if ids1[i] > ids2[i] {
				return 1
			} else if ids1[i] < ids2[i] {
				return -1
			}
		}
	}

	if len(ids1) > len(ids2) {
		return 1
	} else if len(ids1) < len(ids2) {
		return -1
	}

	return 0
}
