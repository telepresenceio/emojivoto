ARG image_registry=ghrc.io/telepresenceio
ARG image_tag=0.1.0
FROM $image_registry/emojivoto-svc-base:$image_tag

ARG svc_name

COPY $svc_name/target/ /usr/local/bin/

# ARG variables arent available for ENTRYPOINT
ENV SVC_NAME $svc_name
ENTRYPOINT cd /usr/local/bin && $SVC_NAME
