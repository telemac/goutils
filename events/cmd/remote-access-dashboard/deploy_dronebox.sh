GOOS=linux go build -o remote-access-dashboard && \
rsync -avzP remote-access-dashboard plugis@dash.idronebox.com:/data/cloud/
