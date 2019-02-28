# Dockerfile for pet-store sample

FROM node:11.2.0-slim
LABEL maintainer="dev@cellery.io"

RUN mkdir -p /app/catalog
WORKDIR /app/catalog
COPY package*.json ./
COPY app.js .

RUN mkdir -p /app/catalog/resources
COPY resources/data.json resources/

RUN npm install
CMD node app.js
EXPOSE 80
