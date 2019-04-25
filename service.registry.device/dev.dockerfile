FROM python:3.6

# Upgrade pip and install YAPF
RUN pip install --upgrade pip && \
    pip install yapf

# Add the libraries
COPY ./libraries/python /root/.local/lib/python3.5/site-packages/

# Install requirements
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY ./service.registry.device/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application files
COPY . .

ARG FLASK_ENV="development"
ENV FLASK_ENV="${FLASK_ENV}" \
    PYTHONUNBUFFERED="true"

ENV FLASK_APP="device_registry"
ENV FLASK_RUN_HOST="0.0.0.0"
ENV FLASK_RUN_PORT="80"

#ENTRYPOINT ["python"]
CMD ["flask", "run"]