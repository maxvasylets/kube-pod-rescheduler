FROM golang:1.12-alpine3.9


COPY build/kube-pod-rescheduler /bin/kube-pod-rescheduler
RUN chmod +x /bin/kube-pod-rescheduler

CMD ["/bin/kube-pod-rescheduler", "--help"]
