FROM node:latest
RUN useradd -ms /bin/bash user_technoservs

WORKDIR home/user_technoservs/

COPY package.json .
RUN npm install --quiet
COPY . ./
EXPOSE 3000