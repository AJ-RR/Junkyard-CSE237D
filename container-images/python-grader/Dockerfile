FROM python:3.13-slim
COPY harness.py /app/harness.py
COPY tests /app/tests
RUN apt-get update && apt-get install -y unzip && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN pip install "gradescope-utils>=0.3.1"
CMD ["python", "/harness.py"]
