FROM busybox

COPY custom-go /fb/custom-go
COPY exported /fb/exported
COPY generated-sdk /fb/generated-sdk
COPY store /fb/store
COPY upload /fb/upload
COPY template /fb/template

CMD [ "cp", "-rfT", "/fb", "/app" ]