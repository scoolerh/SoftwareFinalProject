FROM python:3.11-alpine

WORKDIR /code

COPY . .

RUN pip install flask
RUN pip install requests

ENV FLASK_APP=api.py
ENV FLASK_RUN_HOST=0.0.0.0

CMD ["flask", "run"]

EXPOSE 9000