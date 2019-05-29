datasets.list <- function() {
    script.dir = getSrcDirectory(function(x) {x})
    files.datasets = list.files(path=script.dir, pattern="[.][rR]$")
    files.datasets = files.datasets[!files.datasets %in% c("datasets.R")]
    files.datasets = lapply(files.datasets, tools::file_path_sans_ext, USE.NAMES=FALSE)
    print(files.datasets)
}
