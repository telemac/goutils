GOOS=linux go build -o ansible && \
rsync -avP ansible site.yml plugis@debian:
