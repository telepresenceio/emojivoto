FROM alpine

ARG svc_name

COPY $svc_name/target/ /usr/local/bin/

# ARG variables arent available for ENTRYPOINT
ENV SVC_NAME=$svc_name
ENTRYPOINT [ "sh", "-c", "/usr/local/bin/$SVC_NAME" ]
