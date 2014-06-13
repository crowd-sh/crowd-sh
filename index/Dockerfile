FROM    ubuntu

RUN         apt-get update
RUN         apt-get -y install redis-server

RUN mkdir -p /conf
RUN echo "daemonize yes" > /conf/redis.conf

# Bundle app source
ADD . /code

# Install app dependencies
RUN cd /code

EXPOSE 6379
EXPOSE 3000
CMD ["sh", "/code/run.sh"]