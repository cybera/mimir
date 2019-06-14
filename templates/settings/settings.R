settings.script.dir = dirname(rprojroot::thisfile())
settings.file = file.path(settings.script.dir, "{{.ProjectSettingsPath}}")

settings = RcppTOML::parseTOML(settings.file)
