FROM scratch
ADD .configuration/configuration.json /c.json
ADD appbin /
CMD ["/appbin"]