FROM continuumio/anaconda3:2019.03

RUN pip install toml
RUN /opt/conda/bin/conda install -y r-essentials fonts-anaconda jupyter
RUN Rscript -e "install.packages(c('rprojroot', 'RcppTOML'), repos='http://cran.us.r-project.org')"
