FROM golang

RUN apt update && apt install -y python python-pip
WORKDIR /tmp
COPY requirements.txt /tmp
RUN pip install -r requirements.txt