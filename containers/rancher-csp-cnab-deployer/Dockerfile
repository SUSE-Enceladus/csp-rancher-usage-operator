FROM registry.opensuse.org/opensuse/helm:latest

COPY app/charts /cnab/app/charts
COPY app/run /cnab/app/run
COPY Dockerfile cnab/Dockerfile

ENTRYPOINT [ "/cnab/app/run" ]
