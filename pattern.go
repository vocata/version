package version

import "regexp"

var (
	// pep-508: https://peps.python.org/pep-0508/#names
	packageNameRe = regexp.MustCompile(`(?i)^([A-Z0-9]|[A-Z0-9][A-Z0-9._-]*[A-Z0-9])$`)

	// pip: https://github.com/pypa/pip/blob/23.0.1/src/pip/_internal/models/wheel.py#L15
	wheelFilenameRe = regexp.MustCompile(`(?i)^(?P<namever>(?P<name>[^\s-]+?)-(?P<ver>[^\s-]*?))((-(?P<build>\d[^-]*?))?-(?P<pyver>[^\s-]+?)-(?P<abi>[^\s-]+?)-(?P<plat>[^\s-]+?)\.whl|\.dist-info)$`)

	// pip: https://github.com/pypa/pip/blob/23.0.1/src/pip/_internal/index/package_finder.py#L114
	pyVersionMatchRe = regexp.MustCompile(`-py([123]\.?[0-9]?)$`)

	// pep-440: https://peps.python.org/pep-0440/#appendix-b-parsing-version-strings-with-regular-expressions
	versionRe               = regexp.MustCompile(`^([1-9][0-9]*!)?(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*))*((a|b|rc)(0|[1-9][0-9]*))?(\.post(0|[1-9][0-9]*))?(\.dev(0|[1-9][0-9]*))?$`)
	versionMatchRe          = regexp.MustCompile(`([1-9][0-9]*!)?(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*))*((a|b|rc)(0|[1-9][0-9]*))?(\.post(0|[1-9][0-9]*))?(\.dev(0|[1-9][0-9]*))?`)
	irregularVersionRe      = regexp.MustCompile(`(?i)^\s*v?(?:(?:(?P<epoch>[0-9]+)!)?(?P<release>[0-9]+(?:\.[0-9]+)*)(?P<pre>[-_.]?(?P<pre_l>(a|b|c|rc|alpha|beta|pre|preview))[-_.]?(?P<pre_n>[0-9]+)?)?(?P<post>(?:-(?P<post_n1>[0-9]+))|(?:[-_.]?(?P<post_l>post|rev|r)[-_.]?(?P<post_n2>[0-9]+)?))?(?P<dev>[-_.]?(?P<dev_l>dev)[-_.]?(?P<dev_n>[0-9]+)?)?)(?:\+(?P<local>[a-z0-9]+(?:[-_.][a-z0-9]+)*))?\s*$`)
	irregularVersionMatchRe = regexp.MustCompile(`(?i)v?(?:(?:(?P<epoch>[0-9]+)!)?(?P<release>[0-9]+(?:\.[0-9]+)*)(?P<pre>[-_.]?(?P<pre_l>(a|b|c|rc|alpha|beta|pre|preview))[-_.]?(?P<pre_n>[0-9]+)?)?(?P<post>(?:-(?P<post_n1>[0-9]+))|(?:[-_.]?(?P<post_l>post|rev|r)[-_.]?(?P<post_n2>[0-9]+)?))?(?P<dev>[-_.]?(?P<dev_l>dev)[-_.]?(?P<dev_n>[0-9]+)?)?)(?:\+(?P<local>[a-z0-9]+(?:[-_.][a-z0-9]+)*))?`)
)
