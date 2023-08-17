package version

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Package struct {
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

func NewPackage(name string) (*Package, error) {
	name = CanonicalizePackage(name)
	if !packageNameRe.MatchString(name) {
		return nil, fmt.Errorf("illegal package name '%s'", name)
	}

	return &Package{
		name: name,
	}, nil
}

func (p *Package) String() string {
	return fmt.Sprintf("Package<%s>", p.name)
}

func (p *Package) Name() string {
	return p.name
}

// EvaluateVersion extracts version from filename of current package, original implementations can
// refer to https://github.com/pypa/pip/blob/23.0.1/src/pip/_internal/index/package_finder.py#L108.
func (p *Package) EvaluateVersion(filename string) (string, error) {
	fragment, ext := splitFilename(filename)

	if ext = strings.ToLower(ext); StandardExt.Contains(ext) { // keep the same as pip
		var version string
		if ext == ExtWhl {
			whl, err := NewWheel(filename)
			if err != nil {
				return "", err
			}
			if CanonicalizePackage(whl.Name) != p.name {
				return "", fmt.Errorf("package name '%s' doesn't match that in filename '%s'", p.name, filename)
			}
			version = whl.Version
		}

		if version == "" {
			if version = p.extractVersionFromFragment(fragment); version == "" {
				return "", fmt.Errorf("version not found in '%s'", filename)
			}
		}

		if scope := pyVersionMatchRe.FindStringIndex(version); scope != nil {
			version = version[:scope[0]]
		}

		v, err := Parse(version)
		if err != nil {
			return "", fmt.Errorf("illegal version '%s'", version)
		}

		return v.Complete(), nil // return detailed version
	} else if LegacyExt.Contains(ext) { // rules are casual and used to extract as many versions as possible from filenames
		version := p.extractVersionFromLegacyFragment(fragment)

		if version == "" {
			return "", fmt.Errorf("version not found in '%s'", filename)
		}

		v, err := Parse(version)
		if err != nil {
			return "", fmt.Errorf("illegal version '%s'", version)
		}

		return v.Base(), nil // return brief version
	} else {
		return "", fmt.Errorf("unsupported file extension '%s'", ext)
	}
}

// extractVersionFromFragment extracts version from fragment of source distribution, the filename
// of 'sdist' defined in pep625. For more details: https://peps.python.org/pep-0625/. In addition,
// this method can also extract version from filename of obsolete "bdist_dumb" described in pep527.
// The 'bdist_dumb' will produce files named something like package-1.0.macosx-10.11-x86_64.tar.gz,
// and with the legacy pre PEP 440 versions, 1.0-macosx-10.11-x86_64 is a legal version, for more
// detail, see https://peps.python.org/pep-0527/#bdist-dumb and test cases.
func (p *Package) extractVersionFromFragment(fragment string) string {
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
func (p *Package) extractVersionFromLegacyFragment(fragment string) string {
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
		return nil, fmt.Errorf("illegal wheel filename '%s'", filename)
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
