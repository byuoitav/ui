FROM gcr.io/distroless/static
MAINTAINER Daniel Randall <danny_randall@byu.edu>

COPY ui-linux-amd64 ui
COPY dragonfruit dragonfruit

CMD ["/ui"]
