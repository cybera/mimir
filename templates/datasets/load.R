{{.Name}}.load <- function() {
    script.dir = getSrcDirectory(function(x) {x})
    destfile = file.path(script.dir, "{{.RelPath}}")

    data = data.frame()
    if (file.exists(destfile)) {
        data = read.csv(file=destfile)
    } else {
        stop("Dataset missing from disk")
    }

    return(data)
}
