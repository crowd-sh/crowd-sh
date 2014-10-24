FROM ubuntu:latest

# Bundle app source
ADD . /code

# Install app dependencies
RUN cd /code

EXPOSE 6379
EXPOSE 3000
CMD ["sh", "/code/run.sh"]
