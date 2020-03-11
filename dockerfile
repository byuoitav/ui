FROM gcr.io/distroless/static
MAINTAINER Daniel Randall <danny_randall@byu.edu>

ARG NAME

COPY ${NAME} /ui
COPY dragonfruit /dragonfruit
COPY blueberry /blueberry
COPY cherry /cherry

ENTRYPOINT ["/ui"]
