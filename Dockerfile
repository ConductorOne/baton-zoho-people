FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-zoho-people"]
COPY baton-zoho-people /