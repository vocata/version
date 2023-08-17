package version

import (
	"testing"
)

func TestValidPackageName(t *testing.T) {
	var validPackageNames = []string{
		"0-core-clien",
		"01d61084-d29e-11e9-96d1-7c5cf84ffe8e",
		"19cs30055-q2",
		"a-pandas-ex-columns-and-inde",
		"a2ml",
		"abaqus2dyna",
		"kratoscsharpwrapperapplication",
		"odoo12-addon-beesdoo-stock-coverag",
		"pyqode3-cor",
		"rayaq001",
		"youtube-livechat-scraper-ohn0",
		"nucliadb-utils",
	}

	for _, name := range validPackageNames {
		if _, err := NewPackage(name); err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidPackage(t *testing.T) {
	var invalidPackageNames = []string{
		"mavis-config-",
		"rayane,",
		"rayaq^001",
		"-0",
		"raylib-1.0.",
		"maybe type",
	}

	for _, name := range invalidPackageNames {
		if _, err := NewPackage(name); err == nil {
			t.Error(err)
		}
	}
}

func TestEvaluateVersion(t *testing.T) {
	var packageCases = []struct {
		pkg      string
		filename string
		expected string
		failed   bool
	}{
		{"django-saladoplayer", "django-saladoplayer-0.3.6.tgz", "0.3.6", false},
		{"quicktester", "quicktester-0.3.tgz", "0.3", false},
		{"mother", "mother-0.5.3-r1.tgz", "0.5.3.post1", false},
		{"termfeed", "TermFeed-0.tgz", "0", false},
		{"redis-sniffer", "redis-sniffer_1.0.0.tgz", "", true},
		{"texttable", "texttable-0.8.1.tgz", "0.8.1", false},
		{"generator-tools", "generator_tools-0.3.2-py2.5.tgz", "0.3.2", false},
		{"swiftly", "swiftly-2.00.tgz", "2.0", false},
		{"swiftly", "swiftly-2.04.tgz", "2.4", false},
		{"dipy", "dipy-0.6.0-py2.7-macosx10.6.dmg", "0.6.0", false},
		{"python-igraph", "python_igraph-0.7.1_1-py2.6-macosx10.9.dmg", "0.7.1", false},
		{"fiximports", "fiximports-0.0.1.dev2-py2.py3-none-any.whl", "0.0.1.dev2", false},
		{"fiximports", "fiximports-0.1.15-py2.py3-none-any.whl", "0.1.15", false},
		{"embo", "embo-0.4.0-py3-none-any.whl", "0.4.0", false},
		{"bt-python-sdk", "bt_python_sdk-0.1.0-py3-none-any.whl", "0.1.0", false},
		{"nupyprop", "nupyprop-0.1.6a0.post34-cp39-cp39-macosx_10_9_x86_64.whl", "0.1.6a0.post34", false},
		{"nupyprop", "nupyprop-0.1.7-cp38-cp38-manylinux_2_17_x86_64.manylinux2014_x86_64.whl", "0.1.7", false},
		{"nupyprop", "nupyprop-0.1.7.post1-cp39-cp39-macosx_10_9_x86_64.whl", "0.1.7.post1", false},
		{"nupyprop", "nupyprop-0.1.7.post22-cp39-cp39-manylinux_2_17_x86_64.manylinux2014_x86_64.whl", "0.1.7.post22", false},
		{"nupyprop", "nupyprop-0.1.7.post94-cp38-cp38-macosx_10_9_x86_64.whl", "0.1.7.post94", false},
		{"cnvrg", "cnvrg-0.0.3.5-py3-none-any.whl", "0.0.3.5", false},
		{"cnvrg", "cnvrg-0.1.7-py3-none-any.whl", "0.1.7", false},
		{"simcx", "simcx-1.0.0b16-py2.py3-none-any.whl", "1.0.0b16", false},
		{"sqlite-ulid", "sqlite_ulid-0.2.1a2-py3-none-win_amd64.whl", "0.2.1a2", false},
		{"sqlite-ulid", "sqlite_ulid-0.2.1a9-py3-none-macosx_11_0_arm64.whl", "0.2.1a9", false},
		{"sqlite-ulid", "sqlite_ulid-0.2.1a12-py3-none-macosx_10_6_x86_64.whl", "0.2.1a12", false},
		{"pptx-to-html-lukeehassel", "pptx_to_html_lukeehassel-0.0.1-py3-none-any.whl", "0.0.1", false},
		{"wikiext", "wikiext-0.0.9-py3-none-any.whl", "0.0.9", false},
		{"teletvg-karjakak", "teletvg_karjakak-1.3.0rc1-py3-none-any.whl", "1.3.0rc1", false},
		{"odoo12-addon-product-lot-sequence", "odoo12_addon_product_lot_sequence-12.0.1.0.1.99.dev4-py3-none-any.whl", "12.0.1.0.1.99.dev4", false},
		{"odoo12-addon-product-lot-sequence", "odoo12_addon_product_lot_sequence-12.0.1.0.2-py3-none-any.whl", "12.0.1.0.2", false},
		{"samdata-terminal-dev", "samdata_terminal_dev-2.5.0-py3-none-any.whl", "2.5.0", false},
		{"odoo9-addon-crm-deduplicate-by-website", "odoo9_addon_crm_deduplicate_by_website-9.0.1.0.0.99.dev4-py2-none-any.whl", "9.0.1.0.0.99.dev4", false},
		{"types-aiobotocore-dynamodbstreams", "types_aiobotocore_dynamodbstreams-2.4.1-py3-none-any.whl", "2.4.1", false},
		{"odoo14-addon-l10n-ar-partner", "odoo14_addon_l10n_ar_partner-14.0.0.0.2-py3-none-any.whl", "14.0.0.0.2", false},
		{"dissect-vmfs", "dissect.vmfs-3.2.dev5-py3-none-any.whl", "3.2.dev5", false},
		{"odoo12-addon-fieldservice-maintenance", "odoo12_addon_fieldservice_maintenance-12.0.1.0.0.99.dev1-py3-none-any.whl", "12.0.1.0.0.99.dev1", false},
		{"ausbills", "ausbills-0.4.0-py3-none-any.whl", "0.4.0", false},
		{"socialname", "socialname-0.1.3-py3-none-any.whl", "0.1.3", false},
		{"odoo13-addon-stock-storage-type-putaway-abc", "odoo13_addon_stock_storage_type_putaway_abc-13.0.2.1.0-py3-none-any.whl", "13.0.2.1.0", false},
		{"pynorm", "pynorm-0.1.0-py2.py3-none-any.whl", "0.1.0", false},
		{"pynorm", "pynorm-0.1.5-py2.py3-none-any.whl", "0.1.5", false},
		{"rkd", "rkd-0.1.1.dev16-py3-none-any.whl", "0.1.1.dev16", false},
		{"rkd", "rkd-0.6.0.0b2.dev6-py3-none-any.whl", "0.6.0.0b2.dev6", false},
		{"tencentcloud-sdk-python-autoscaling", "tencentcloud_sdk_python_autoscaling-3.0.388-py2.py3-none-any.whl", "3.0.388", false},
		{"shore-kafka", "shore-kafka-2022.1.8.tar.gz", "2022.1.8", false},
		{"tptp-lark-parser", "tptp_lark_parser-0.1.2.tar.gz", "0.1.2", false},
		{"python-dhl", "python-dhl-1.0.0.dev6.tar.gz", "1.0.0.dev6", false},
		{"python-dhl", "python-dhl-1.0.0.dev12.macosx-10.6-x86_64.tar.gz", "1.0.0.dev12.macosx-10.6-x86_64", false}, // refer to https://peps.python.org/pep-0527/#bdist-dumb
		{"infi-dtypes-hctl", "infi.dtypes.hctl-0.0.4-develop-1-g07c2bdb.tar.gz", "0.0.4-develop-1-g07c2bdb", false},
		{"twython", "twython-1.2.macosx-10.5-i386.tar.gz", "1.2.macosx-10.5-i386", false},
		{"twython", "twython-1.2.tar.gz", "1.2", false},
		{"scikit-kinematics", "scikit-kinematics-0.8.3.win-amd64.exe", "0.8.3", false},
		{"esclient", "ESClient-0.5.4.macosx-10.8-intel.exe", "0.5.4", false},
		{"sloth-ci-ext-docker-exec", "sloth-ci.ext.docker_exec-1.0.8.win-amd64.exe", "1.0.8", false},
		{"pyteomics-biolccc", "pyteomics.biolccc-1.5.0.win-amd64-py2.6.exe", "1.5.0", false},
		{"pyxer", "pyxer-0.6.1.win32.exe", "0.6.1", false},
		{"epyunit", "epyunit-0.2.8.linux-x86_64.exe", "0.2.8", false},
		{"configviper", "ConfigViper-0.1.win32-py2.6.exe", "0.1", false},
		{"python-streamtools", "python-streamtools-0.0.4.macosx-10.9-intel.exe", "0.0.4", false},
		{"gemfire-rest", "gemfire-rest-1.0.macosx-10.9-intel.exe", "1.0", false},
		{"goldsaxcreatetablesyfinance", "GoldSaxCreateTablesYFinance-1.01.win-amd64.exe", "1.1", false},
		{"coal-mine", "coal_mine-0.4-1.noarch.rpm", "0.4", false},
		{"ll-core", "ll-core-1.9.1-1.i386.rpm", "1.9.1", false},
		{"polib", "polib-0.3.0-1.noarch.rpm", "0.3.0", false},
		{"epyunit", "epyunit-0.1.9-1.noarch.rpm", "0.1.9", false},
		{"python-foreman", "python-foreman-0.2.1-1.src.rpm", "0.2.1", false},
		{"turboflot", "python-turboflot-0.0.9-1.fc8.src.rpm", "", true},
		{"turboflot", "python-turboflot-0.1.1-1.fc9.noarch.rpm", "", true},
		{"lmdb", "lmdb-0.81-py3.4-win-amd64.egg", "0.81", false},
		{"lmdb", "lmdb-0.89-py3.4-win32.egg", "0.89", false},
		{"lmdb", "lmdb-0.97-py3.6-win-amd64.egg", "0.97", false},
		{"lmdb", "lmdb-1.0.0-py2.7-win32.egg", "1.0.0", false},
		{"suddendeath", "suddendeath-0.1.0-py2.7.egg", "0.1.0", false},
		{"pythondata-cpu-ibex", "pythondata_cpu_ibex-0.0.post2061-py2.7.egg", "0.0", false},
		{"pythondata-cpu-ibex", "pythondata_cpu_ibex-0.0.post2061-py3.8.egg", "0.0", false},
		{"postgpd", "postgpd-1.2.zip", "1.2", false},
		{"gecosistema-lite", "gecosistema_lite-0.0.617.zip", "0.0.617", false},
		{"dugong", "dugong-3.0.tar.bz2", "3.0", false},
		{"pypol2", "pypol_-0.5.tar.bz2", "", true},
		{"fqueue", "fqueue-0.0.5.tar.bz2", "0.0.5", false},
		{"collective-portlet-paypal", "collective.portlet.paypal-0.2dev-r57320.tar.bz2", "0.2dev-r57320", false},
		{"django-selectel-storage", "django-selectel-storage-0.3.1.tar.bz2", "0.3.1", false},
		{"appwsgi", "657.tar.bz2", "", true},
		{"appwsgi", "appwsgi 667.tar.bz2", "", true},
		{"appwsgi", "appwsgi 1014.tar.bz2", "", true},
		{"pyarmor", "pyarmor-2.0.1.tar.bz2", "2.0.1", false},
		{"winappdbg", "winappdbg-1.3.win32-py2.6.msi", "1.3", false},
		{"lorm", "lorm-0.2.11.win32.msi", "0.2.11", false},
		{"experimentdb", "experimentdb-0.2.win32.msi", "0.2", false},
		{"pytango", "PyTango-8.0.2.win32-py3.2.msi", "8.0.2", false},
		{"pycegui", "PyCEGUI-0.7.5.win32-py2.6.msi", "0.7.5", false},
		{"py-postgresql", "py-postgresql-1.0.0.win32-py3.1.msi", "1.0.0", false},
		{"garlicsim-wx", "garlicsim_wx-0.4.win32-py2.6.msi", "0.4", false},
		{"argparse", "argparse-1.0.1.win32.msi", "1.0.1", false},
		{"pygresql", "PyGreSQL-4.2.1.win32-py2.6.msi", "4.2.1", false},
		{"uselesscapitalquiz", "uselesscapitalquiz-3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196442881097566593-py3-none-any.whl", "3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196442881097566593", false},
		{"uselesscapitalquiz", "uselesscapitalquiz-3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196442881097566593.tar.gz", "3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709384460955058223172535940812848111745028410270193852110555964462294895493038196442881097566593", false},
		{"12-test", "12@test-0.1.tar.gz", "", true},
	}

	for _, c := range packageCases {
		t.Run(c.pkg, func(t *testing.T) {
			pkg, err := NewPackage(c.pkg)
			if err != nil {
				t.Error(err)
				return
			}

			version, err := pkg.EvaluateVersion(c.filename)
			if c.failed != (err != nil) {
				t.Error(err)
				return
			}

			if version != c.expected {
				t.Errorf("actual: %s, %s != %s", version, version, c.expected)
			}
		})
	}
}
