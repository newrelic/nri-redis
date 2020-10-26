FROM alpine

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
RUN echo "AWS_DEFAULT_REGION: $AWS_DEFAULT_REGION"

CMD ["/entrypoint.sh"]