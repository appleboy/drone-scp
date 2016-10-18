FROM centurylink/ca-certs

ADD drone-sftp /

ENTRYPOINT ["/drone-sftp"]
