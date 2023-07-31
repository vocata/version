package version

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var ErrUndefinedVersion = errors.New("undefined version")

type PyPiPackage struct {
	name string
}

// CanonicalizePackage standardizes a package name, for detail:
// https://peps.python.org/pep-0503/#normalized-names
var irregularPackageLetters = regexp.MustCompile(`[-_.]+`)

func CanonicalizePackage(name string) string {
	return strings.ToLower(irregularPackageLetters.ReplaceAllString(name, "-"))
}

var irregularLegacyPackageLetters = regexp.MustCompile(`[^0-9a-zA-Z]+`)

func CanonicalizeLegacyPackage(name string) string {
	return strings.ToLower(irregularLegacyPackageLetters.ReplaceAllString(name, "-"))
}

func NewPackage(name string) (*PyPiPackage, error) {
	name = CanonicalizePackage(name)
	if !packageNameRe.MatchString(name) {
		return nil, fmt.Errorf("invalid package name '%s'", name)
	}

	return &PyPiPackage{
		name: name,
	}, nil
}

func (p *PyPiPackage) String() string {
	return p.name
}

// EvaluateVersion extracts version from filename of current package, original implementations can
// refer to https://github.com/pypa/pip/blob/23.0.1/src/pip/_internal/index/package_finder.py#L108.
func (p *PyPiPackage) EvaluateVersion(filename string) (string, error) {
	fragment, ext := splitFilename(filename)

	if ext = strings.ToLower(ext); standardExt.Contains(ext) { // keep the same as pip
		var version string
		if ext == ExtWhl {
			whl, err := NewWheel(filename)
			if err != nil {
				return "", fmt.Errorf("%w, %s", ErrUndefinedVersion, err.Error())
			}
			if CanonicalizePackage(whl.Name) != p.name {
				return "", fmt.Errorf("%w, package name in '%s' not matched", ErrUndefinedVersion, filename)
			}
			version = whl.Version
		}

		if version == "" {
			if version = p.extractVersionFromFragment(fragment); version == "" {
				return "", fmt.Errorf("%w, version not found in '%s'", ErrUndefinedVersion, filename)
			}
		}

		if scope := pyVersionMatchRe.FindStringIndex(version); scope != nil {
			version = version[:scope[0]]
		}

		v, err := Parse(version)
		if err != nil {
			return "", fmt.Errorf("%w, not a valid version '%s'", ErrUndefinedVersion, version)
		}

		return v.Complete(), nil // return detailed version
	} else if legacyExt.Contains(ext) { // rules are casual and used to extract as many versions as possible from filenames
		version := p.extractVersionFromLegacyFragment(fragment)

		if version == "" {
			return "", fmt.Errorf("%w, version not found in '%s'", ErrUndefinedVersion, filename)
		}

		v, err := Parse(version)
		if err != nil {
			return "", fmt.Errorf("%w, not a valid version '%s'", ErrUndefinedVersion, version)
		}

		return v.Base(), nil // return brief version
	} else {
		return "", fmt.Errorf("%w, unsupported extension '%s'", ErrUndefinedVersion, ext)
	}
}

// extractVersionFromFragment extracts version from fragment of source distribution, the filename
// of 'sdist' defined in pep625. For more details: https://peps.python.org/pep-0625/. In addition,
// this method can also extract version from filename of obsolete "bdist_dumb" described in pep527.
// The 'bdist_dumb' will produce files named something like package-1.0.macosx-10.11-x86_64.tar.gz,
// and with the legacy pre PEP 440 versions, 1.0-macosx-10.11-x86_64 is a valid, for more detail,
// see https://peps.python.org/pep-0527/#bdist-dumb and test cases.
func (p *PyPiPackage) extractVersionFromFragment(fragment string) string {
	for i, c := range fragment {
		if c != '-' {
			continue
		}
		if CanonicalizePackage(fragment[:i]) == p.name {
			return fragment[i+1:]
		}
	}

	return ""
}

// extractVersionFromLegacyFragment extracts version from fragment of legacy distribution such as
// rpm, deb, exe etc. Since there is no officially defined specification, this method tries the
// best to extract the version number from the fragment using loose rules.
func (p *PyPiPackage) extractVersionFromLegacyFragment(fragment string) string {
	for i, c := range fragment {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			continue
		}
		if CanonicalizeLegacyPackage(fragment[:i]) == p.name {
			return irregularVersionMatchRe.FindString(fragment[i+1:])
		}
	}

	return ""
}

type Wheel struct {
	Filename string
	Name     string
	Version  string
	Btag     string
	Pyvers   []string
	Abis     []string
	Plats    []string
}

// NewWheel create a Wheel object from filename, which contains five segments of filename
// described in https://peps.python.org/pep-0427/#file-name-convention, implementation refers
// to https://github.com/pypa/pip/blob/23.0.1/src/pip/_internal/models/wheel.py#L22.
func NewWheel(filename string) (*Wheel, error) {
	match := wheelFilenameRe.FindStringSubmatch(filename)
	if match == nil {
		return nil, fmt.Errorf("'%s' is not a valid wheel filename", filename)
	}

	return &Wheel{
		Filename: filename,
		Name:     strings.ReplaceAll(match[wheelFilenameRe.SubexpIndex("name")], "_", "-"),
		Version:  strings.ReplaceAll(match[wheelFilenameRe.SubexpIndex("ver")], "_", "-"),
		Btag:     match[wheelFilenameRe.SubexpIndex("build")],
		Pyvers:   strings.Split(match[wheelFilenameRe.SubexpIndex("pyver")], "."),
		Abis:     strings.Split(match[wheelFilenameRe.SubexpIndex("abi")], "."),
		Plats:    strings.Split(match[wheelFilenameRe.SubexpIndex("plat")], "."),
	}, nil
}
