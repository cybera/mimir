package languages

var Supported = []string{
	"python",
	"r",
}

var Extensions = map[string]string{
	"python": ".py",
	"r":      ".R",
}

var InitFiles = map[string]map[string]string{
	"python": {
		"datasets/__init__.py": "src/datasets/__init__.py",
	},
	"r": {
		"datasets/datasets.R": "src/datasets/datasets.R",
	},
}
