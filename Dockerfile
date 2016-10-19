FROM centurylink/ca-certs

ADD drone-scp /

ENTRYPOINT ["/drone-scp"]
