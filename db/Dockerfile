FROM python:3.11-alpine

WORKDIR /code

COPY . .

RUN pip install flask

ENV FLASK_APP=data.py
ENV FLASK_RUN_HOST=0.0.0.0

CMD ["flask", "run"]

EXPOSE 7788