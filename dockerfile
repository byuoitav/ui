FROM gcr.io/distroless/static
MAINTAINER Daniel Randall <danny_randall@byu.edu>

COPY ui-linux-x86 ui
COPY dragonfruit-dist dragonfruit

CMD ["/ui"]
