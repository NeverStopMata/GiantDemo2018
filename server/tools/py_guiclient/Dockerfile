FROM python

RUN git clone https://github.com/openx/python3-protobuf.git
RUN cd python3-protobuf && git clone https://github.com/google/googletest.git && mv googletest gtest && cd gtest && git checkout release-1.5.0
RUN cd python3-protobuf && ./autogen.sh && ./configure && make && cd python && python setup.py install

ADD dockerhelp/sources.list /etc/apt/sources.list
RUN apt-get update
RUN yes | apt-get install libgtk-3-dev
RUN yes | apt-get install build-essential
RUN yes | apt-get install freeglut3-dev
RUN yes | apt-get install libgstreamer0.10-dev
RUN yes | apt-get install libgstreamer-plugins-base0.10-dev
RUN yes | apt-get install libgstreamer0.10
RUN yes | apt-get install libgstreamer-plugins-base0.10
RUN yes | apt-get install libwebkitgtk-3.0-dev
RUN pip install wxPython

RUN mkdir -p /py_guiclient/proto
RUN mkdir -p /py_guiclient/gamesync
RUN mkdir -p /conf
COPY *.py /py_guiclient/
COPY proto/*.py /py_guiclient/proto/
COPY gamesync/*.py /py_guiclient/gamesync/

VOLUME /conf

CMD [ "python", "/py_guiclient/_pystress_main.py", "--cfg=/conf/cfg.json" ]
