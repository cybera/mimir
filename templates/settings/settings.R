settings.script.dir = getSrcDirectory(function(x) {x})
settings.file = file.path(settings.script.dir, "{{.ProjectSettingsPath}}")

settings = RcppTOML::parseTOML(settings.file)
