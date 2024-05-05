FROM python:3.8.13
LABEL authors="sefatascan"

WORKDIR /app

COPY . .

RUN pip install --upgrade pip
RUN pip install -r ./requirements.txt

CMD ["python", "main.py"]