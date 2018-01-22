FROM python

ADD sources.list /etc/apt/sources.list
RUN apt-get update
RUN yes | apt-get install curl

COPY *.py /

ENTRYPOINT [ "python", "/crontab.py" ]
